import requests


class RatelimitValidator():
    def __init__(self, shell, gateway):
        self.shell = shell
        self.gateway = gateway
    
    def validate(self, domain, path, retry, ratelimited):
        if self.gateway:
            print("validate using port forward")
            headers = {
                'Host': domain,
            }
            
            for sequence in range(retry):
                response = requests.get('http://localhost:8080%s' % path, headers=headers)
                if ratelimited:
                    if response.status_code != 429:
                        if sequence is not retry-1:
                            continue
                        raise Exception("response code: %d, it's not ratelimited" % response.status_code) 
                else:
                    if response.status_code == 429:
                        if sequence is not retry-1:
                            continue
                        raise Exception("response code: %d, it's ratelimited" % response.status_code) 
            
            print("validate using port kubectl exec")
            validate_command = ["kubectl", "-n", "development", "exec", "-i", "deploy/client", "-c", "client",
                                 "--", "curl", "http://istio-ingressgateway.istio-system.svc.cluster.local:80%s" %(path), "-H", "'Host:", "%s'" %(domain), "--write-out", "'%{json}'"]
            
            for sequence in range(retry):
                out = self.shell.os_execute(' '.join(validate_command))
                if ratelimited:
                    if '"http_code":429' not in out:
                        if sequence is not retry-1:
                            continue   
                        raise Exception("it's not ratelimited")
                else:
                    if '"http_code":429' in out:
                        if sequence is not retry-1:
                            continue   
                        raise Exception("it's ratelimited")          
        
        else:
            print("validate using port kubectl exec")
            validate_command = ["kubectl", "-n", "development", "exec", "-i", "deploy/client", "-c", "client",
                              "--", "curl", "http://%s:9898%s" %(domain, path), "-H", "'Host:", "%s'" %(domain), "--write-out", "'%{json}'"]
            for sequence in range(retry):
                out = self.shell.os_execute(' '.join(validate_command))
                if ratelimited:
                    if '"http_code":429' not in out:
                        if sequence is not retry-1:
                            continue
                        raise Exception("it's not ratelimited") 
                else:
                    if '"http_code":429' in out:
                        if sequence is not retry-1:
                            continue
                        raise Exception("it's not ratelimited")               
            
        print('istio-ratelimit-operator is working!')

