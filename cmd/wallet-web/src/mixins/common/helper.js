/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import jp from 'jsonpath';
import { v4 as uuidv4 } from 'uuid';
const { encodeURI } = require('js-base64');
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

function contextCacheClosure() {
  let contextCache = {};

  return function (contextURI) {
    let cacheElement = contextCache[contextURI];
    if (cacheElement !== undefined) {
      return cacheElement;
    }

    return new Promise((resolve, reject) => {
      let req = new XMLHttpRequest();

      req.onload = (e) => {
        if (req.status !== 200) {
          console.log('fetching remote context failed with status code ' + req.status);
          reject('failed with status code ' + req.status);
        }

        cacheElement = JSON.parse(req.responseText);
        contextCache[contextURI] = cacheElement;
        resolve(cacheElement);
      };

      req.open('GET', contextURI);
      req.send();
    });
  };
}

export const getContext = contextCacheClosure();

export function matchTypeInContext(ctxObj, type) {
  return getTermInContext(ctxObj, type);
}

// splits the iri if it's a compact IRI, returns undefined otherwise.
function getCompactIRISplit(iri) {
  let idx = iri.indexOf(':');
  if (idx === -1) return undefined;

  let pref = iri.substring(0, idx);
  if (pref === 'http' || pref === 'https') return undefined; // assume http and https prefixes are literal.

  let suff = iri.substring(idx + 1);

  return [pref, suff];
}

// Get a term in a context, expanding compact IRIs as necessary
function getTermInContext(ctxObj, term) {
  let res;

  if (Array.isArray(ctxObj)) {
    for (let i = 0; i < ctxObj.length; i++) {
      let curr = ctxObj[i];

      // string in a context array is a remote context reference.
      // note: dereferencing isn't needed for bank account demo
      // which is all that uses the legacy js-side presexch filtering.
      if (typeof curr === 'string' || curr instanceof String) {
        continue;
      }

      res = getTermInContext(ctxObj[i], term);
      if (!!res) {
        break;
      }
    }
  } else if (typeof ctxObj === 'string' || ctxObj instanceof String) {
    // if the parent object contained a string at index [term], then we've found an instance of term.
    // in this case, we return immediately, so the parent can handle any IRI expansion.
    return ctxObj;
  } else if (typeof ctxObj === 'object' && !!ctxObj) {
    // if this is an object (we've already handled arrays above), check if it contains term
    // either directly, or within a @context member.
    res = getTermInContext(ctxObj[term], term);
    if (!res) {
      res = getTermInContext(ctxObj['@context'], term);
    }
    if (!res) {
      res = getTermInContext(ctxObj['@id'], term); // @id holds the ID of a json-ld type
    }
  }

  // term not present inside ctx
  if (!res) {
    return undefined;
  }

  // res is an IRI, if it's compact we need to try and expand it.
  let compactSplit = getCompactIRISplit(res);
  if (!compactSplit || !compactSplit.length || compactSplit.length < 2) {
    // res is not a compact IRI, so we return it directly.
    return res;
  }

  // res is a compact IRI, so we try to lookup its prefix within the current ctxObj.
  let pref = getTermInContext(ctxObj, compactSplit[0]);
  if (!pref) {
    // res is a compact IRI, but it must be expanded in a parent context, not the current one.
    return res;
  }

  // concatenate expanded prefix with original suffix
  return pref + compactSplit[1];
}

// function to get the credential display data
export function getCredentialDisplayData(vc, manifest, skipEmpty = true) {
  const id = encodeURI(populatePath(vc, manifest.id));
  const brandColor = manifest.brandColor || '';
  const issuanceDate = populatePath(vc, manifest.issuanceDate);
  const title = populatePath(vc, manifest.title.path) || manifest.title.fallback;
  const icon = manifest.icon;

  // find properties
  const _readProperty = (property) => {
    const value = populatePath(vc, property.path) || '';
    const label = property.label;
    const type = property.type;
    const format = property.format;
    return {
      label,
      value,
      type,
      format,
    };
  };

  // fallback credential will use vc credentialSubject as properties
  const properties =
    Object.values(manifest.properties).length === 0
      ? Object.entries(vc.credentialSubject).map(([label, value]) => ({
          label,
          value,
        }))
      : Object.values(manifest.properties)
          .map(_readProperty)
          .filter((el) => !(skipEmpty && el.value.length === 0));

  return {
    id,
    brandColor,
    issuanceDate,
    title,
    icon,
    properties,
  };
}
// Populates path property in the JSON object
function populatePath(vc, paths) {
  for (const path of paths) {
    try {
      const resolvedQuery = jp.value(vc, path);
      if (resolvedQuery) return resolvedQuery;
    } catch (error) {
      // TODO: write this error into logger once implemented (it would mean we received corrupt value from config file)
      console.warn('failed to read display data from credential config', error);
    }
  }
  return undefined;
}

export const toLower = (text) => text.toString().toLowerCase();

export const minsToNanoSeconds = (ns) => ns * 60 * 10 ** 9;

export const getDIDVerificationMethod = (dids, id) => {
  return jp.query(dids, `$[?(@.id=="${id}")].verificationMethod[*].id`);
};

export function wait(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export function getCredentialIcon(staticAssetsUrl, icon) {
  if (staticAssetsUrl) {
    return `${staticAssetsUrl}/images/icons/${icon}`;
  }
  return `${require('@/assets/img/credential--generic-icon.svg')}`;
}

/**
 *  Finds matching credential manifest output descriptors from given list of manifest mappings.
 *
 *  If found, prepares fresh credential manifest from output descriptors found.
 *  If not found, creates fresh credential manifest using all fields found in credential subject of credential.
 *
 */
export function prepareCredentialManifest(presentation, manifestDictionary, issuer) {
  const _findOutputDescriptor = (ctxmap, credential) => {
    if (!ctxmap) {
      return;
    }

    for (let ctx of credential['@context']) {
      if (ctxmap[ctx]) {
        return ctxmap[ctx];
      }
    }
  };

  let credentialManifest = {
    id: uuidv4(),
    issuer: {
      id: issuer || uuidv4(),
    },
    version: '0.1.0',
    output_descriptors: [],
  };
  let uniqueOutputDescriptors = new Set();

  let fulfillment = {
    id: uuidv4(),
    manifest_id: credentialManifest.id,
    descriptor_map: [],
  };

  presentation.verifiableCredential.forEach((credential, index) => {
    // find output descriptor for given credential type.
    let entry = _findOutputDescriptor(
      manifestDictionary[getCredentialType(credential.type)],
      credential
    );
    if (entry) {
      uniqueOutputDescriptors.add(...entry['output_descriptors']);
    } else {
      // find default output descriptor
      entry = _findOutputDescriptor(manifestDictionary['VerifiableCredential'], credential);

      if (entry) {
        // populate output descriptor properties
        for (let k of Object.keys(credential.credentialSubject)) {
          entry.output_descriptors[0].display.properties.push({
            path: [`$.credentialSubject.${k}`],
            schema: {
              type: 'string',
            },
            label: k,
          });
        }

        uniqueOutputDescriptors.add(...entry['output_descriptors']);
      } else {
        console.error(
          "couldn't find default credential manifest for given credential",
          credential['@context'],
          credential['type']
        );
        throw "couldn't find default credential manifest for given credential";
      }
    }

    // prepare fulfillment,
    fulfillment.descriptor_map.push({
      id: entry['output_descriptors'][0].id,
      format: 'ldp_vc',
      path: `$.verifiableCredential[${index}]`,
    });
  });

  presentation['credential_fulfillment'] = fulfillment;

  credentialManifest.output_descriptors = Array.from(uniqueOutputDescriptors);
  return credentialManifest;
}

/**
 * Check credential manifest for missing optional properties before resolving to ensure the
 * resolved manifest has all fields required for wallet to display the credential properly
 *
 *  @param {Object} credentialManager - Credential Manager to resolve manifest.
 *  @param {Object} manifestDictionary - Default credential manifest output descriptors.
 *  @param {String} auth - authorization token for wallet operations.
 *  @param {Object} options - options to resolve credential from wallet.
 *  @param {String} options.credentialID - (optional) ID of the credential to be resolved from wallet content store.
 *  @param {String} options.credential - (optional) raw credential data model to be resolved.
 *  @param {String} options.fulfillment - (optional) credential fulfillment using which given raw credential or credential ID to be resolved.
 *  @param {String} options.manifestID - (optional) ID of the manifest from wallet content store.
 *  @param {String} options.manifest - (optional) raw manifest to be used for resolving credential.
 *  @param {String} options.descriptorID - (optional) if fulfillment not provided then this descriptor ID can be used to resolve credential.
 *
 *  @returns {Promise<Object>} - promise containing resolved results or error if operation fails.
 */
export function resolveManifest(
  credentialManager,
  manifestDictionary,
  auth,
  { credentialID, credential, fulfillment, manifestID, manifest, descriptorID }
) {
  const cred = credential || fulfillment.verifiableCredential[0];
  const filledManifest = checkManifestForMissingFields(
    manifest,
    getCredentialType(cred.type),
    manifestDictionary
  );
  return credentialManager.resolveManifest(auth, {
    credentialID,
    credential,
    fulfillment,
    manifestID,
    manifest: filledManifest,
    descriptorID,
  });
}

function checkManifestForMissingFields(manifest, credentialType, manifestDictionary) {
  const outputDescriptor = manifest.output_descriptors[0];
  const defaultOutputDescriptor = getDefaultOutputDescriptor(credentialType, manifestDictionary);
  if (!outputDescriptor.display) {
    outputDescriptor.display = defaultOutputDescriptor.display;
  } else {
    const displayProperties = ['title', 'subtitle', 'description'];
    displayProperties.forEach((property) => {
      if (isEmpty(outputDescriptor.display[property])) {
        outputDescriptor.display[property] = defaultOutputDescriptor.display[property];
      }
    });
  }
  if (!outputDescriptor.styles) {
    outputDescriptor.styles = defaultOutputDescriptor.styles;
  } else {
    const styleProperties = ['thumbnail', 'hero', 'background', 'text'];
    styleProperties.forEach((property) => {
      if (isEmpty(outputDescriptor.styles[property])) {
        outputDescriptor.styles[property] = defaultOutputDescriptor.styles[property];
      }
    });
  }
  manifest.output_descriptors[0] = outputDescriptor;

  return manifest;
}

function getDefaultOutputDescriptor(credentialType, manifestDictionary) {
  let defaultManifest;
  if (!Object.prototype.hasOwnProperty(manifestDictionary, credentialType)) {
    console.warn('default credential manifest for credential type does not exist');
    defaultManifest = manifestDictionary['VerifiableCredential'];
  } else {
    defaultManifest = manifestDictionary[credentialType];
  }
  const outputDescriptor = defaultManifest[Object.keys(defaultManifest)[0]].output_descriptors[0];
  return outputDescriptor;
}

// Check if the object or any of its sub-objects are empty
function isEmpty(property) {
  if (!property) return true;
  for (const key in property) {
    if (Object.keys(property[key]).length === 0) return true;
  }
  return false;
}
