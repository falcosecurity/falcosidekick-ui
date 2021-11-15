const crypto = require('crypto');

export const utils = {
  priorityToColor(priority) {
    switch (priority) {
      case 'Emergency':
        return '#C62828';
      case 'Alert':
        return '#D32F2F';
      case 'Critical':
        return '#E53935';
      case 'Error':
        return '#FF5252';
      case 'Warning':
        return '#FB8C00';
      case 'Notice':
        return '#1976D2';
      case 'Informational':
        return '#03A9F4';
      case 'Debug':
        return '#29B6F6';
      default:
        return '#555';
    }
  },
  stringToColor(str) {
    return `#${crypto.createHash('md5').update(str).digest('hex').substring(0, 6)}`;
  },
};

export default {
};
