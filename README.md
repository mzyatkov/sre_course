# sre_course

### Run

#### Priklad
```bash
helm upgrade --install priklad-chart priklad-chart --kubeconfig ./student1.yaml -f ./priklad-chart/customvalues.yaml 
```

#### Monitoring
```bash
cd monitoring
./download_debs.sh
ansible-playbook -i inventory.ini site.yml
```