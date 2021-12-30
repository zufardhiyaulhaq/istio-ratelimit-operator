import requests


class RatelimitValidator():
    def __init__(self, shell, gateway):
        self.shell = shell
        self.gateway = gateway
    
    def validate(self, domain, path, retry):
        if self.gateway:
            headers = {
                'Host': domain,
            }

            for sequence in range(retry):
                response = requests.get('http://localhost:8080%s' % path, headers=headers)
                if response.status_code != 429:
                    if sequence is not retry:
                        continue
                    raise Exception("response code: %d, it's not ratelimited" % response.status_code) 
            
            validate_command = ["kubectl", "-n", "development", "exec", "-i", "deploy/client", "-c", "client",
                                 "--", "curl", "http://istio-ingressgateway.istio-system.svc.cluster.local:80%s" %(path), "-H", "'Host:", "%s'" %(domain), "--write-out", "'%{json}'"]
            
            for sequence in range(retry):
                out = self.shell.os_execute(' '.join(validate_command))
                print(out)
                if '"http_code":429' not in out:
                    if sequence is not retry:
                        continue   
                    raise Exception("it's not ratelimited") 
        
        else:
            validate_command = ["kubectl", "-n", "development", "exec", "-i", "deploy/client", "-c", "client",
                              "--", "curl", "http://%s:9898%s" %(domain, path), "-H", "'Host:", "%s'" %(domain), "--write-out", "'%{json}'"]
            for sequence in range(retry):
                out = self.shell.os_execute(' '.join(validate_command))
                print(out)
                if '"http_code":429' not in out:
                    if sequence is not retry:
                        continue
                    raise Exception("it's not ratelimited") 
            
        print('ratelimit is working!')

