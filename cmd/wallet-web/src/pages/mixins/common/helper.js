/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import jp from 'jsonpath';
import { PresentationExchange } from './presentationExchange';

var flatten = require('flat');

const ALL_ICONS = [
  'account_box',
  'contacts',
  'person',
  'person_outline',
  'card_membership',
  'portrait',
  'bento',
  'directions_car',
  'house',
];
const VC_ICONS_MAP = {
  PermanentResidentCard: 'house',
  UniversityDegreeCredential: 'school',
  BookingReferenceCredential: 'flight',
  VaccinationCertificate: 'health_and_safety',
  mDL: 'directions_car',
};

// flatten credential subjects
export function flatCredentialSubject(subj) {
  return flatten(subj, {
    transformKey: function (key) {
      let parts = key.split('#');
      return parts[parts.length - 1];
    },
  });
}

// extracts query reasons from the list of queries.
export function extractQueryReasons(query) {
  return jp.query(query, `$[*].credentialQuery..reason`);
}

// extract all query types form the list of queries.
export function extractQueryTypes(query) {
  return jp.query(query, `$[*].type`);
}

// extracts presentation exchange reasons from presentation exchange query found in query list.
export function extractPresentationExchangeReasons(query) {
  let pexch = jp.query(query, `$[?(@.type=="PresentationExchange")]`);

  let reqs = pexch.map((p) => {
    if (p.credentialQuery.length > 1) {
      throw 'unsupported operation, can have only one presentation exchange inside credential query';
    }

    return new PresentationExchange(p.credentialQuery[0]).requirementDetails();
  });

  return reqs.reduce((acc, val) => acc.concat(...val), []);
}

// normalizeQuery fixes if credential query format as not per this wallet's standards.
export function normalizeQuery(query) {
  let _normalize = (q) => {
    q.credentialQuery = Array.isArray(q.credentialQuery) ? q.credentialQuery : [q.credentialQuery];
    return q;
  };

  let allQuery = Array.isArray(query) ? query : [query];
  return allQuery.map(_normalize);
}

export const filterCredentialsByType = (creds, types, include) =>
  creds.filter((c) =>
    include ? types.includes(getCredentialType(c.type)) : !types.includes(getCredentialType(c.type))
  );

export const getVCIcon = (type) =>
  VC_ICONS_MAP[type]
    ? VC_ICONS_MAP[type]
    : ALL_ICONS[Math.floor(Math.random() * Math.floor(ALL_ICONS.length))];

export const isVPType = (type) => toLower(type) == 'verifiablepresentation';

export const getCredentialType = (types) =>
  types.filter((type) => type != 'VerifiableCredential')[0];

export const toLower = (text) => text.toString().toLowerCase();

export const minsToNanoSeconds = (ns) => ns * 60 * 10 ** 9;

export const getDIDVerificationMethod = (dids, id) => {
  return jp.query(dids, `$[?(@.id=="${id}")].verificationMethod[*].id`);
};
