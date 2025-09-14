const { InfluxDB, Point } = require("@influxdata/influxdb-client");

const url = process.env.INFLUX_URL || "http://localhost:8086";
const token = process.env.INFLUX_TOKEN || 'energy-token';
const org = process.env.INFLUX_ORG || 'energy-org';
const bucket = process.env.INFLUX_BUCKET || 'energy-bucket';

const client = new InfluxDB({ url, token });
const writeApi = client.getWriteApi(org, bucket, 'ns');

function writeMeasurement(data) {
  try {
    const point = new Point('energy')
      .tag('user_id', data.user_id)
      .tag('user_type', data.user_type)
      .floatField('consumption', data.energy_consumption)
      .floatField('generation', data.energy_generation)
      .floatField('storage', data.energy_storage)
      .stringField('simulation_time', data.timestamp)
      .timestamp(new Date());

    console.log(data.energy_generation)

    writeApi.writePoint(point);
    console.log(`[INFLUX] Gravado: ${data.user_id} em ${data.timestamp}`);
    
  } catch (err) {
    console.error('[INFLUX] Erro ao gravar:', err);
  }
}

module.exports = { writeMeasurement };