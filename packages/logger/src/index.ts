export type LogLevel = "debug" | "info" | "warn" | "error";

export interface LogContext {
  requestId?: string;
  service?: string;
  [key: string]: unknown;
}

export interface Logger {
  debug(message: string, context?: LogContext): void;
  info(message: string, context?: LogContext): void;
  warn(message: string, context?: LogContext): void;
  error(message: string, context?: LogContext): void;
}

function write(level: LogLevel, message: string, context: LogContext = {}) {
  const payload = {
    level,
    message,
    timestamp: new Date().toISOString(),
    ...context
  };

  console[level === "debug" ? "log" : level](JSON.stringify(payload));
}

export function createLogger(defaultContext: LogContext = {}): Logger {
  return {
    debug: (message, context) => write("debug", message, { ...defaultContext, ...context }),
    info: (message, context) => write("info", message, { ...defaultContext, ...context }),
    warn: (message, context) => write("warn", message, { ...defaultContext, ...context }),
    error: (message, context) => write("error", message, { ...defaultContext, ...context })
  };
}

