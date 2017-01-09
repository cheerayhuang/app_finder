import argparse
import math
import requests
import sys
import subprocess
import time
import types

try:
    import xml.etree.cElementTree as ET
except ImportError:
    import xml.etree.ElementTree as ET

URL_APPLE = 'http://localhost:5001/apple/'
URL_GOOGLE = 'http://localhost:5001/google/'
URL_NOT_FOUND = 'http://localhost:5001/notfound/'


DAY = 24*3600
CMD = 'psql -h sdkbox.crzhtaa6feqd.us-west-2.redshift.amazonaws.com -p 5439 "user=awsuser password=678D9ed29db6 dbname=events" < /tmp/tmp.sql'

SQL = '''
\o /tmp/out.txt;
select distinct sdkbox_app_package_id
from event_log
where timestamp > extract(epoch from '2016-12-13'::date) and timestamp < extract(epoch from '2016-12-14'::date)
limit 10;
'''

def print_log(s):
    print time.strftime('[%Y-%m-%d %H:%M:%S]', time.localtime()) + s
    sys.stdout.flush()


def get_time():
    global SQL
    start = time.strftime('%Y-%m-%d', time.localtime(time.time()-2*DAY))
    end = time.strftime('%Y-%m-%d', time.localtime(time.time()-DAY))

    SQL = SQL.replace('2016-12-13', start)
    SQL = SQL.replace('2016-12-14', end)
    SQL = SQL.replace('out.txt', 'out-' + start + '.txt')

    in_f = open('/tmp/tmp.sql', 'w')
    in_f.write(SQL)
    in_f.close()

    return start, end


def refind_not_found():
    print_log('find bundleId from Notfound table first...')
    sys.stdout.flush()
    r = requests.post(URL_NOT_FOUND)
    l = r.json()

    for k in l:
        try:
            r = requests.get(URL_APPLE + k, timeout=10)
        except Exception:
            time.sleep(1)
            continue

        r_map = r.json()
        if 'err' in r_map or r_map[k] == "not found":
            try:
                r = requests.get(URL_GOOGLE + k, timeout=10)
            except Exception:
                continue

        r_map = r.json()
        if 'err' in r_map or r_map[k] == "not found":
            print_log(r.text)
            sys.stdout.flush()
            continue

        r = requests.delete(URL_NOT_FOUND + k)
        print_log(r.text)
        sys.stdout.flush()


def find(start, end):
    f = '/tmp/out-' + start + '.txt'
    print_log('execute query, generate result into {0}...'.format(f))
    '''
    r = subporcess.call(CMD)
    if r != 0:
        print_log('execute query failed')
        sys.exit(1)

    '''
    print_log('find bundleId from file: {0}...'.format(f))
    in_file = open(f, 'r')
    lines = in_file.readlines()
    in_file.close()

    # trim "filed name" and "(x) rows".
    lines = lines[2:-1]

    total_bundleid = 0
    total_not_found = 0
    total_err = 0
    for text in lines:
        not_access_api = False
        t1 = time.time()
        if type(text) is types.StringType and not text.isdigit():
            text = text.lower()
            total_bundleid += 1
            try:
                r = requests.get(URL_APPLE + text, timeout=10)
            except Exception:
                time.sleep(1)
                continue
            r_map = r.json()
            if 'not_access_apple_api' in r_map:
                not_access_api = True
            if 'err' in r_map or r_map[text] == "not found":
                try:
                    r = requests.get(URL_GOOGLE + text, timeout=10)
                except Exception:
                    continue

            if text in r.json() and r.json()[text] == "not found":
                total_not_found += 1
                r = requests.post(URL_NOT_FOUND + text)

            if 'err' in r.json():
                total_err += 1

            print_log(r.text)

            diff = time.time() - t1 - 3;
            #if diff < 0 and not not_access_api:
            #    time.sleep(math.ceil(-diff))

    print_log('total: ' + str(total_bundleid))
    print_log('total "not found": ' + str(total_not_found))
    print_log('total "errors": ' + str(total_err))


def main():
    '''
    project_name = 'client'

    parser = argparse.ArgumentParser(prog=project_name)
    parser.add_argument('-x', '--xml', default='./list.xml',help="specify a bunldID list in xml format.", required=True)

    args = parser.parse_args()
    '''

    start, end = get_time()
    refind_not_found()
    time.sleep(5)
    find(start, end)

if __name__ == '__main__':

    main()
