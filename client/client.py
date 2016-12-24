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

def find(args):
    ns = {'ss': 'urn:schemas-microsoft-com:office:spreadsheet'}
    tree = ET.ElementTree(file=args.xml)

    d = tree.find('ss:Worksheet/ss:Table/ss:Row/ss:Cell/ss:Data', ns)
    for elem in tree.iterfind('ss:Worksheet/ss:Table/ss:Row', ns):
        for data in elem.iterfind('ss:Cell/ss:Data', ns):
            text = data.text
            not_access_api = False
            t1 = time.time()
            if type(text) is types.StringType and not text.isdigit():
                text = text.lower()
                try:
                    r = requests.get(URL_APPLE + text, timeout=10)
                except requests.exceptions.Timeout:
                    continue
                r_map = r.json()
                if 'not_access_apple_api' in r_map:
                    not_access_api = True
                if 'err' in r_map or r_map[text] == "not found":
                    try:
                        r = requests.get(URL_GOOGLE + text, timeout=10)
                    except requests.exceptions.Timeout:
                        continue
                print r.text
                diff = time.time() - t1 - 3;
                if diff < 0 and not not_access_api:
                    time.sleep(math.ceil(-diff))



def main():
    project_name = 'client'

    parser = argparse.ArgumentParser(prog=project_name)
    parser.add_argument('-x', '--xml', default='./list.xml',help="specify a bunldID list in xml format.", required=True)

    args = parser.parse_args()

    find(args)

if __name__ == '__main__':

    main()
