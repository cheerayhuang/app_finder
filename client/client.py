import argparse
import requests
import sys
import time
import types

try:
    import xml.etree.cElementTree as ET
except ImportError:
    import xml.etree.ElementTree as ET

URL = 'http://test.sdkbox:5001/apple/'

def find(args):
    ns = {'ss': 'urn:schemas-microsoft-com:office:spreadsheet'}
    tree = ET.ElementTree(file=args.xml)

    d = tree.find('ss:Worksheet/ss:Table/ss:Row/ss:Cell/ss:Data', ns)
    for elem in tree.iterfind('ss:Worksheet/ss:Table/ss:Row', ns):
        for data in elem.iterfind('ss:Cell/ss:Data', ns):
            text = data.text
            if type(text) is types.StringType and not text.isdigit():
                try:
                    r = requests.get(URL + text, timeout=10)
                except Exception as e:
                    continue
                r_map = r.json()
                print r_map
                if 'http_count' in r_map or 'err' in r_map:
                    time.sleep(3)

def main():
    project_name = 'client'

    parser = argparse.ArgumentParser(prog=project_name)
    parser.add_argument('-x', '--xml', default='./list.xml',help="specify a bunldID list in xml format.", required=True)

    args = parser.parse_args()

    find(args)

if __name__ == '__main__':

    main()
