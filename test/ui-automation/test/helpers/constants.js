/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

'use strict';

const constants = {};
module.exports = constants;

constants.dids = {
    trustbloc: {
        name: 'DID_TRUSTBLOC',
        keyType: 'Ed25519',
        signatureType: 'Ed25519Signature2018',
    },
    key: {
        name: 'DID_KEY',
        did: 'did:key:z6MknbvvxApNNLmhZV8JZXsoXcppnaF2dgG4bqN5Xi74Pisq',
        pkjwk: '{ "kty": "OKP", "d": "q6-Z0JFBQZRBoFwyAcpWolHZPWgmvi3VcyPXggwO2lo", "crv": "Ed25519", "x": "eRYcIeukCnDTrRXa6qgYoQpTfmcZqZzoDMznqDir7-g", "kid": "z6MknbvvxApNNLmhZV8JZXsoXcppnaF2dgG4bqN5Xi74Pisq"}',
        keyID: 'did:key:z6MknbvvxApNNLmhZV8JZXsoXcppnaF2dgG4bqN5Xi74Pisq#z6MknbvvxApNNLmhZV8JZXsoXcppnaF2dgG4bqN5Xi74Pisq',
        signatureType: 'Ed25519Signature2018'
    },
    v1: {
        name: 'DID_V1',
        did: 'did:v1:test:nym:z6Mkg1G7NEzfu7Wc937oeQuZjerX4jM1kvBMWQWX3oFoUvTR',
        pkjwk: '{"kty": "OKP","d":"Dq5t2WS3OMzcpkh8AyVxJs5r9v4L39ocIz9CpUOqM40","crv": "Ed25519","x": "ODaPFurJgFcoVCUYEmgOJpWOtPlOYbHSugasscKWqDM","kid":"z6MkiEh8RQL83nkPo8ehDeE7cPHPEkA5dDUeFHbtJJF8Sc2v"}',
        keyID: 'did:v1:test:nym:z6Mkg1G7NEzfu7Wc937oeQuZjerX4jM1kvBMWQWX3oFoUvTR#z6MkiEh8RQL83nkPo8ehDeE7cPHPEkA5dDUeFHbtJJF8Sc2v',
        signatureType: 'Ed25519Signature2018'
    }
};
