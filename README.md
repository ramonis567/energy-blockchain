If grafana not work:
sudo chown -R 472:472 ./grafana-data


sequência para iniciar o sistema:

Atualizar go para 1.21 -> em /blockchain   chmod +x update-go.sh   -> ./update-go.sh

em /blockchain: ./start.sh
em /simulation -> start.py (terminal dedicado)
em /interconnection -> start.py (terminal dedicado)


portas: 
http://localhost:8086   - Influx db    (admin/admin123)
http://localhost:3001   - Grafana      (admin/admin123)
http://localhost:1883   - MQTT   

- backend
- frontend