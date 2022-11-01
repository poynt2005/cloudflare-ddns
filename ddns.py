import urllib.request
from venv import create

from CloudFlare import CloudFlare
import argparse, os
from pprint import pprint

CLOUDFLARE_DNS_API_KEY = os.environ.get('CLOUDFLARE_DNS_API_KEY', None)
CLOUDFLARE_DNS_ZONE_NAME = os.environ.get('CLOUDFLARE_DNS_ZONE_NAME', None)
CLOUDFLARE_DNS_SUBDOMAIN_NAME = os.environ.get('CLOUDFLARE_DNS_SUBDOMAIN_NAME', None)

def get_external_ip():
    return urllib.request.urlopen('https://ident.me').read().decode('utf-8')

def renew_dns():
    cf = CloudFlare(token=CLOUDFLARE_DNS_API_KEY)

    zone_id = ''
    try:
        params = { 'name': CLOUDFLARE_DNS_ZONE_NAME }
        zones = cf.zones.get(params=params)
        zone_id = zones[0]['id']
    except Exception as e:
        print('Error happened: ' + str(e))
        print('/zones.get api call failed')
        return False

    find_target_subdomain = None

    try:
        dns_records = cf.zones.dns_records.get(zone_id)
        a_records = []
        
        for raw_record in dns_records:
            if raw_record['type'] == 'A':
                a_records.append({
                    'ip': raw_record['content'],
                    'domain': raw_record['name'],
                    'id': raw_record['id']
                })

        for record in a_records:
            if record['domain'].split('.')[0] == CLOUDFLARE_DNS_SUBDOMAIN_NAME:
                find_target_subdomain = record['id']
                break
                
    except Exception as e:
        print('Error happened: ' + str(e))
        print('/zones/dns_records/export api call failed')
        return False

    if find_target_subdomain is None:
        print('Create DNS record: %s.%s' % (CLOUDFLARE_DNS_SUBDOMAIN_NAME, CLOUDFLARE_DNS_ZONE_NAME))
        try:
            create_info = cf.zones.dns_records.post(zone_id, data={
                'type': 'A',
                'name': CLOUDFLARE_DNS_SUBDOMAIN_NAME,
                'content': get_external_ip()
            })
        except Exception as e:
            print('Error happened: ' + str(e))
            print('/zones/dns_records/create api call failed')
            return False
    else:
        print('Update DNS record: %s.%s' % (CLOUDFLARE_DNS_SUBDOMAIN_NAME, CLOUDFLARE_DNS_ZONE_NAME))

        try:
            put_info = cf.zones.dns_records.put(zone_id, find_target_subdomain, data={
                'type': 'A',
                'name': CLOUDFLARE_DNS_SUBDOMAIN_NAME,
                'content': get_external_ip()
            })
        except Exception as e:
            print('Error happened: ' + str(e))
            print('/zones/dns_records/update api call failed')
            return False

    return True

if __name__ == '__main__':
    argparser = argparse.ArgumentParser()
    argparser.add_argument('-k', '--key', help='Provide your cloudflare DNS api key', type=str)
    argparser.add_argument('-z', '--zone', help='Provide your cloudflare DNS zone', type=str)
    argparser.add_argument('-d', '--domain', help='Provide your cloudflare subdomain name', type=str)

    args = argparser.parse_args()

    if args.key:
        CLOUDFLARE_DNS_API_KEY = args.key
    
    if args.zone:
        CLOUDFLARE_DNS_ZONE_NAME = args.zone

    if args.domain:
        CLOUDFLARE_DNS_SUBDOMAIN_NAME = args.domain
    
    if CLOUDFLARE_DNS_API_KEY is None or CLOUDFLARE_DNS_SUBDOMAIN_NAME is None or CLOUDFLARE_DNS_ZONE_NAME is None:
        print('cloudflare arguments are none, exiting program...')
        exit()

    if renew_dns():
        print('Domain name: %s.%s renewed a new DNS record' % (CLOUDFLARE_DNS_SUBDOMAIN_NAME, CLOUDFLARE_DNS_ZONE_NAME))
    else:
        print('Domain name: %s.%s renewed a new DNS record failed' % (CLOUDFLARE_DNS_SUBDOMAIN_NAME, CLOUDFLARE_DNS_ZONE_NAME))

