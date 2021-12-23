import requests

class RatelimitValidator():
    def __init__(self):
        pass
    
    def validate(self, domain, path):
        headers = {
            'Host': domain,
        }

        response = requests.get('http://localhost:8080%s' % path, headers=headers)
        if response.status_code != 429:
            raise Exception("It's not ratelimited") 
        
        print('Ratelimit works!')
