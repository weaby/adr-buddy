// @decision.id: adr-1
// @decision.name: Using Pino for logging
// @decision.status: accepted
// @decision.category: infrastructure
// @decision.context: We needed structured logging with low overhead
//   for our high-throughput API. After benchmarking several options,
//   Pino provided the best performance characteristics.
// @decision.decision: Adopt Pino as the standard logging library
//   across all Node.js services.
// @decision.consequences: All services must migrate from Winston.
//   We gain 3-5x better logging performance.

const pino = require('pino');
const logger = pino();

logger.info('Application started');
