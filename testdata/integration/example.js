// @decision.id: adr-1
// @decision.name: Using Pino for logging
// @decision.status: accepted
// @decision.category: infrastructure
// @decision.context: We needed structured logging with low overhead
//   for our high-throughput API. After benchmarking Winston, Bunyan,
//   and Pino, we found that Pino provided the best performance
//   characteristics with the lowest memory overhead.
// @decision.decision: Adopt Pino as the standard logging library
//   across all Node.js services. All new services must use Pino,
//   and existing services should migrate during their next major version.
// @decision.consequences: All services must migrate from Winston.
//   We gain 3-5x better logging performance but lose some
//   Winston-specific plugins. Team needs training on Pino's API.

const pino = require('pino');
const logger = pino({
  level: process.env.LOG_LEVEL || 'info'
});

module.exports = logger;
