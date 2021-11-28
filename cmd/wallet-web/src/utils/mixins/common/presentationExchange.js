/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Ajv from 'ajv';
import jp from 'jsonpath';
import { presentationDefSchema } from './presentationDefSchema';
import { getContext, getCredentialType, matchTypeInContext } from '..';

const presentationSubmissionTemplate = `{
    "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://identity.foundation/presentation-exchange/submission/v1"
    ],
    "type": ["VerifiablePresentation", "PresentationSubmission"],
    "presentation_submission": {
        "descriptor_map": []
    },
    "verifiableCredential": []
}`;

const defSubmissionRuleName = 'Requested information';
const defSubmissionRulePurpose = 'We need below information from your wallet';
const defSubmissionRule = 'all conditions should be met';
const manifestCredType = 'IssuerManifestCredential';

/**
 * PresentationDefinition represents Presentation Definitions objects
 * to articulate what proofs an entity requires to make a decision about an interaction with a Subject.
 *
 * RFC: https://identity.foundation/presentation-exchange
 *
 * Note:
 *  when 2 credentials match to same input descriptor id, showing both in presentation submission descriptor map
 *
 * @param presentation definition requirement object
 * @class
 */
export class PresentationExchange {
  constructor(requirement) {
    validateSchema(requirement);

    this.requirementObjs = requirement['submission_requirements'];
    this.descriptors = requirement['input_descriptors'];
    this.applyRules = this.requirementObjs && this.requirementObjs.length > 0;
    if (this.applyRules) this._filterDescriptors();
  }

  _filterDescriptors() {
    // validate groups defined in 'submission_requirements'
    var requiredRules = jp.query(this.requirementObjs, '$..from');
    var availableRules = jp.query(this.descriptors, '$..group[*]');
    if (!requiredRules.every((v) => availableRules.includes(v))) {
      throw [
        { message: "Couldn't find matching group in descriptors for 'submission_requirements'" },
      ];
    }

    // retain descriptors only needed by required rules
    let descriptors = this.descriptors.filter((descriptor) =>
      descriptor.group.some((v) => requiredRules.includes(v))
    );

    let descrsByGroup = new Map();
    requiredRules.forEach(function (rule) {
      descrsByGroup[rule] = descriptors.filter((descriptor) => descriptor.group.includes(rule));
    });

    this.descriptorsByGroup = descrsByGroup;
    this.descriptors = descriptors;
  }

  async createPresentationSubmission(vcs) {
    let manifests = jp.query(vcs, `$[?(@.type.indexOf('${manifestCredType}') != -1)]`);
    let credentials = jp.query(vcs, `$[?(@.type.indexOf('${manifestCredType}') == -1)]`);

    let results = [];
    if (this.applyRules) {
      results = await evaluateByRules(
        credentials,
        manifests,
        this.descriptorsByGroup,
        this.requirementObjs
      );
    } else {
      results = await evaluateAll(credentials, manifests, this.descriptors);
    }

    return prepareSubmission(results);
  }

  requirementDetails() {
    let result = [];

    if (this.applyRules) {
      let descrsByGroup = this.descriptorsByGroup;
      this.requirementObjs.forEach(function (obj, index) {
        let r = {};
        let { name, purpose, from } = obj;

        r.name = name ? name : `${defSubmissionRuleName} #${index + 1}`;
        r.purpose = purpose ? purpose : defSubmissionRulePurpose;
        r.rule = countDetails(obj);
        r.descriptors = [];

        descrsByGroup[from].forEach(function (d) {
          r.descriptors.push(getNameAndPurpose(d));
        });

        result.push(r);
      });
    } else {
      let r = {
        name: defSubmissionRuleName,
        purpose: defSubmissionRulePurpose,
        rule: defSubmissionRule,
        descriptors: [],
      };

      this.descriptors.forEach(function (descriptor) {
        r.descriptors.push(getNameAndPurpose(descriptor));
      });

      result.push(r);
    }

    return result;
  }
}

function getNameAndPurpose(descriptor) {
  let name = jp.query(descriptor, '$.name');
  let purpose = jp.query(descriptor, '$.purpose');
  let constraints = jp.query(descriptor, '$.constraints.fields[*].purpose');

  return {
    name: name.length > 0 ? name[0] : 'Condition details are not provided in request',
    purpose: purpose.length > 0 ? purpose[0] : '',
    constraints,
  };
}

// Schema validator
var ajv = new Ajv({ $data: true });
var validate = ajv.compile(presentationDefSchema);

function validateSchema(data) {
  let valid = validate(data);
  if (!valid) {
    throw validate.errors;
  }
}

// doesCredMatchDescriptor returns true if given credential matches descriptor, false otherwise.
async function doesCredMatchDescriptor(credential, descriptor) {
  if (getCredentialType(credential.type) === 'IssuerManifestCredential') {
    return false;
  }

  // match schema
  let schemas;
  let requiredSchemas = descriptor.schema.filter((s) => s.required);
  let required = !!requiredSchemas.length;
  if (required) {
    schemas = requiredSchemas.map((s) => s.uri);
  } else {
    schemas = descriptor.schema.map((s) => s.uri);
  }

  let contexts = Array.isArray(credential['@context'])
    ? credential['@context']
    : [credential['@context']];
  let types = Array.isArray(credential['type']) ? credential['type'] : [credential['type']];
  let expandedTypes = [];

  let contextObjs = await Promise.all(contexts.map((ctxURI) => getContext(ctxURI)));

  for (let i = 0; i < contextObjs.length; i++) {
    let ctx = contextObjs[i];

    for (let j = 0; j < types.length; j++) {
      let typ = matchTypeInContext(ctx, types[j]);
      if (typ) {
        expandedTypes.push(typ);
      }
    }
  }

  let schemaMatched = required; // with required schemas, check for failure, otherwise, check for success

  for (let i = 0; i < schemas.length; i++) {
    if (expandedTypes.includes(schemas[i])) {
      if (!required) {
        schemaMatched = true;
        break;
      }
    } else {
      if (required) {
        schemaMatched = false;
        break;
      }
    }
  }

  if (!schemaMatched) {
    // schema not matched, skip this credential
    return false;
  }

  if (
    !descriptor.constraints ||
    !descriptor.constraints.fields ||
    descriptor.constraints.fields.length == 0
  ) {
    // if no constraints declared, credential matched !!
    return true;
  }

  // found constraints, apply filter using constraints,
  let filterMatched;
  for (let f in descriptor.constraints.fields) {
    let field = descriptor.constraints.fields[f];
    let valueFound;
    // look for matching value
    for (let p in field.path) {
      valueFound = jp.query(credential, field.path[p]);
      if (!valueFound || valueFound.length > 0) {
        break;
      }
    }

    // no matching path found in given credential
    if (valueFound && valueFound.length == 0) {
      filterMatched = false;
      break;
    }

    // if filter present, then apply filter
    if (field.filter) {
      valueFound = valueFound.filter((v) => ajv.validate(field.filter, v));
      if (valueFound.length == 0) {
        // only if the result is valid, proceed iterating the rest of the fields entries
        filterMatched = false;
        break;
      }
    }

    filterMatched = true;
  }

  return filterMatched;
}

// matchManifest matches if descriptor schema exists in manifest credential contexts list
// TODO: manifests to have credential previews so that complete constraint checks can be run
function matchManifest(manifest, descriptor) {
  let schemas = descriptor.schema
    .map((s) => s.uri)
    .map((uri) => uri.substring(0, uri.indexOf('#')));

  if (descriptor.constraints && descriptor.constraints.fields) {
    descriptor.constraints.fields
      .filter((f) => f.filter)
      .filter((f) => f.filter.const)
      .forEach((f) => {
        schemas.push(f.filter.const);
      });
  }

  return manifest.credentialSubject.contexts.some((v) => schemas.includes(v));
}

// prepareSubmission creates presentation submission for all matched credentials
function prepareSubmission(results) {
  let presentationSubmission = JSON.parse(presentationSubmissionTemplate);

  results.forEach(function (result, index) {
    //TODO add VC only once if it matches 2 conditions
    presentationSubmission.verifiableCredential.push(result.credential);
    presentationSubmission.presentation_submission.descriptor_map.push({
      id: result.id,
      path: `$.verifiableCredential[${index}]`,
    });
  });

  return presentationSubmission;
}

// evaluateAll evaluates credentials based on all input descriptors
async function evaluateAll(credentials, manifests, descriptors) {
  let collectedResults = await Promise.all(
    descriptors.map((descriptor) =>
      (async function () {
        let credMatches = await Promise.all(
          credentials.map((credential) =>
            doesCredMatchDescriptor(credential, descriptor).then(function (res) {
              if (res) {
                return { credential, id: descriptor.id };
              }

              return false;
            })
          )
        );

        let partialResults = credMatches.filter((res) => !!res);
        let matched = partialResults.length > 0;

        // none of the credentials matched, check for manifest credential matches
        if (!matched && manifests) {
          partialResults = manifests
            .filter((cred) => matchManifest(cred, descriptor))
            .map((credential) => ({ credential, id: descriptor.id, manifest: true }));
        }

        return partialResults;
      })()
    )
  );

  return collectedResults.reduce((a, b) => a.concat(b));
}

// evaluateByRules evaluates credentials based on submission rules
async function evaluateByRules(credentials, manifests, descrsByGroup, submissions) {
  let collectedResults = await Promise.allSettled(
    submissions.map((submission) =>
      (async function () {
        let descriptors = descrsByGroup[submission.from];
        let pick = countMatcher(submission, descriptors.length);

        let credMatches = await Promise.allSettled(
          credentials.map((credential) =>
            (async function () {
              let isMatch = await Promise.allSettled(
                descriptors.map((d) => doesCredMatchDescriptor(credential, d))
              );

              let matches = descriptors.filter((_, i) => !!isMatch[i].value);

              if (pick(matches.length)) {
                return matches.map((match) => ({
                  credential: credential,
                  id: match.id,
                }));
              } else {
                return false;
              }
            })()
          )
        );

        let partialResults = credMatches
          .filter((res) => !!res.value)
          .map((res) => res.value)
          .reduce((a, b) => a.concat(b));

        let matched = partialResults.length > 0;

        // none of the credentials matched, check for manifest credential matches
        if (!matched && manifests) {
          partialResults = manifests
            .map((credential) => ({
              credential,
              matches: descriptors.filter((d) => matchManifest(credential, d)),
            }))
            .filter((res) => pick(res.matches.length))
            .map((res) =>
              res.matches.map((match) => ({
                credential: res.credential,
                id: match.id,
                manifest: true,
              }))
            );
        }

        return partialResults;
      })()
    )
  );

  return collectedResults
    .filter((res) => !!res.value)
    .map((res) => res.value)
    .reduce((a, b) => a.concat(b));
}

function countMatcher(submission, descrLen) {
  return function (matchedCount) {
    if (submission.rule == 'all' && matchedCount == descrLen) {
      return true;
    } else if (submission.count && matchedCount >= submission.count) {
      return true;
    } else if (
      submission.max &&
      submission.min &&
      matchedCount <= submission.max &&
      matchedCount >= submission.min
    ) {
      return true;
    } else if (submission.max && matchedCount <= submission.max) {
      return true;
    } else if (submission.min && matchedCount >= submission.min) {
      return true;
    } else {
      return false;
    }
  };
}

function countDetails(submission) {
  if (submission.rule == 'all') {
    return 'all conditions should be met';
  } else if (submission.count) {
    return `at least ${submission.count} condition(s) should be met`;
  } else if (submission.max && submission.min) {
    return `${submission.min} to ${submission.max} conditions should be met`;
  } else if (submission.max) {
    return `at most ${submission.max} conditions should be met`;
  } else if (submission.min) {
    return `at least $submission.{count} condition(s) should be met`;
  } else {
    return '';
  }
}
