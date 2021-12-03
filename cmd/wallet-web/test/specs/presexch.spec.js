/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { matchTypeInContext } from '@/utils/mixins/common/helper.js';
import { getTestData } from '../common';
import { expect } from 'chai';

let contexts = getTestData('contexts/agent-startup-contexts.json').documents;

function getCtx(url) {
  for (const ctx of contexts) {
    if (ctx.url === url) {
      return ctx.content;
    }
  }

  return {};
}

describe('getting credential type from context', function () {
  it('matches in mDL context', function () {
    let ctx = getCtx('https://trustbloc.github.io/context/vc/examples/mdl-v1.jsonld');

    expect(matchTypeInContext(ctx, 'name')).to.equal('http://schema.org/name');
    expect(matchTypeInContext(ctx, 'mDL')).to.equal('https://example.org/examples#mDL');
  });

  it('matches in authorization credential context', function () {
    let ctx = getCtx('https://trustbloc.github.io/context/vc/authorization-credential-v1.jsonld');

    expect(matchTypeInContext(ctx, 'AuthorizationCredential')).to.equal(
      'https://example.org/examples#AuthorizationCredential'
    );
  });

  it('matches in verifiable credential context', function () {
    let ctx = getCtx('https://www.w3.org/2018/credentials/v1');

    expect(matchTypeInContext(ctx, 'VerifiableCredential')).to.equal(
      'https://www.w3.org/2018/credentials#VerifiableCredential'
    );
    expect(matchTypeInContext(ctx, 'VerifiablePresentation')).to.equal(
      'https://www.w3.org/2018/credentials#VerifiablePresentation'
    );
  });

  it('matches in example credential context', function () {
    let ctx = getCtx('https://www.w3.org/2018/credentials/examples/v1');

    expect(matchTypeInContext(ctx, 'UniversityDegreeCredential')).to.equal(
      'https://example.org/examples#UniversityDegreeCredential'
    );
    expect(matchTypeInContext(ctx, 'RelationshipCredential')).to.equal(
      'https://example.org/examples#RelationshipCredential'
    );
  });

  it('matches in citizenship credential context', function () {
    let ctx = getCtx('https://w3id.org/citizenship/v1');

    expect(matchTypeInContext(ctx, 'PermanentResidentCard')).to.equal(
      'https://w3id.org/citizenship#PermanentResidentCard'
    );
    expect(matchTypeInContext(ctx, 'PermanentResident')).to.equal(
      'https://w3id.org/citizenship#PermanentResident'
    );
  });
});
