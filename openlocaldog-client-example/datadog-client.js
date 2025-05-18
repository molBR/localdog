const StatsD = require('hot-shots');

// Connect to Localdog (UDP on port 8125)
const dogstatsd = new StatsD({
  host: 'localhost',
  port: 8125,
  prefix: 'myapp.',
  globalTags: ['env:dev'],
  protocol: 'udp'
});

// Send a counter metric
dogstatsd.increment('users_logged_in', 1, ['feature:signup']);

// Send a gauge metric
dogstatsd.gauge('latency', 120, ['api:login']);

console.log('Metrics sent to Localdog.');