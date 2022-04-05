import { encodeURI } from 'js-base64';

export default class JSONWebToken {
  constructor(payload, algorithm = 'none') {
    this.header = {
      alg: algorithm,
    };
    this.payload = payload;
    this.signature = {};
  }

  getTokenString() {
    const encodedHeader = encodeURI(this.getHeader());
    const encodedPayload = encodeURI(this.getPayload());
    const encodedSignature = encodeURI(this.getSignature());
    return `${encodedHeader}.${encodedPayload}.${encodedSignature}`;
  }

  getSignature() {
    // TODO: Sign with an appropiate algorithm
    return JSON.stringify(this.signature);
  }

  getHeader() {
    return JSON.stringify(this.header);
  }

  getPayload() {
    return JSON.stringify(this.payload);
  }
}
