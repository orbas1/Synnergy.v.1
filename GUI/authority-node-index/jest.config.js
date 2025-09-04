module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  roots: ['<rootDir>/src', '<rootDir>/tests'],
  collectCoverage: true,
  coverageThreshold: {
    global: {
      lines: 60,
      branches: 50,
    },
  },
};
