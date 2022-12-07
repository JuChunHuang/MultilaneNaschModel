from pywebio import start_server
from pywebio.input import *
from pywebio.output import *
from pywebio.session import info as session_info

def main():
    info = input_group('Input parameters', [
        input("roadLength", name="roadLength", type=NUMBER),
        input("NSDV probability", name="NSDVprobability", type=FLOAT),
    ])

    put_markdown('Your roadLength: `%d`, Category: `%.1f`' % (info['roadLength'], info['NSDVprobability']))
    img = open('D:/cmu_S1/02601/go/src/MultilaneNaschModel/SingleLane.out.gif', 'rb').read()
    put_image(img,width='100px',height='100px')


if __name__ == '__main__':
    start_server(main, debug=True, port=8080)