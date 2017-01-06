import argparse
import math
import requests
import sys
import time
import types

try:
    import xml.etree.cElementTree as ET
except ImportError:
    import xml.etree.ElementTree as ET

URL_APPLE = 'http://localhost:5001/apple/'
URL_GOOGLE = 'http://localhost:5001/google/'
URL_NOT_FOUND = 'http://localhost:5001/notfound/'

def refind_not_found():
    print 'find bundleId from Notfound table first...'
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
            print r.text
            continue

        r = requests.delete(URL_APPLE + k)
        print r.text


def find(args):
    print 'find bundleId from xml file: {0}...'.format(args.xml)
    ns = {'ss': 'urn:schemas-microsoft-com:office:spreadsheet'}
    tree = ET.ElementTree(file=args.xml)

    total_bundleid = 0
    total_not_found = 0
    total_err = 0
    d = tree.find('ss:Worksheet/ss:Table/ss:Row/ss:Cell/ss:Data', ns)
    for elem in tree.iterfind('ss:Worksheet/ss:Table/ss:Row', ns):
        for data in elem.iterfind('ss:Cell/ss:Data', ns):
            text = data.text
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

                print r.text

                diff = time.time() - t1 - 3;
                #if diff < 0 and not not_access_api:
                #    time.sleep(math.ceil(-diff))

    print 'total: ' + str(total_bundleid)
    print 'total "not found": ' + str(total_not_found)
    print 'total "errors": ' + str(total_err)



def main():
    project_name = 'client'

    parser = argparse.ArgumentParser(prog=project_name)
    parser.add_argument('-x', '--xml', default='./list.xml',help="specify a bunldID list in xml format.", required=True)

    args = parser.parse_args()

    refind_not_found()
    time.sleep(5)
    find(args)

if __name__ == '__main__':

    main()
