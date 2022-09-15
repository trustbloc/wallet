/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

'use strict';

import { chapi, wallet } from '../helpers';

const credential = new Map();
credential.set('PermanentResidentCard', {
  name: 'Permanent Resident Card',
  vc: {
    '@context': [
      'https://www.w3.org/2018/credentials/v1',
      'https://w3id.org/citizenship/v1',
      'https://w3id.org/vc-revocation-list-2020/v1',
      'https://w3id.org/security/bbs/v1',
    ],
    credentialStatus: {
      id: 'https://issuer-vcs.trustbloc.local/status/1#3',
      revocationListCredential: 'https://issuer-vcs.trustbloc.local/status/1',
      revocationListIndex: '3',
      type: 'RevocationList2020Status',
    },
    credentialSubject: {
      birthCountry: 'Bahamas',
      birthDate: '1958-07-17',
      familyName: 'Pasteur',
      gender: 'Male',
      givenName: 'Louis',
      id: 'did:orb:AiMP4:EiA2Gtl-qHKjTouzu-Rd0cYOwQxJ-qN0DO0HNnhfXXCqCg',
      lprCategory: 'C09',
      lprNumber: '999-999-999',
      residentSince: '2015-01-01',
      type: ['Person', 'PermanentResident'],
    },
    description: 'Permanent Resident Card of Mr.Louis Pasteur',
    id: 'http://example.com/eb299a34-529e-4a84-a67c-573865db4aa7',
    issuanceDate: '2021-03-11T14:52:00.8492482Z',
    issuer: {
      id: 'did:key:zUC72c7u4BYVmfYinDceXkNAwzPEyuEE23kUmJDjLy8495KH3pjLwFhae1Fww9qxxRdLnS2VNNwni6W3KbYZKsicDtiNNEp76fYWR6HCD8jAz6ihwmLRjcHH6kB294Xfg1SL1qQ',
      name: 'didkey-bbsblssignature2020-bls12381g2',
    },
    name: 'Permanent Resident Card',
    proof: {
      created: '2021-07-06T18:16:57.739627-04:00',
      proofPurpose: 'assertionMethod',
      proofValue:
        'koyKGr8WwjCOUqm-HV_7SVtvIIM4EhnJJ_8P2k0RF3ElQP2ntQJMKKtpoQTqk5l3QI5jN0Zn8nHJm3gyFkKdYJpC4IseNTU98u9UTijHlABpAhGbDaKTHs-b1IDsHkx_DrR3BSktz1Va_cilRP2WqA',
      type: 'BbsBlsSignature2020',
      verificationMethod:
        'did:key:zUC73A7EHiDAxxy29qox4hD5Dyc6fXqStkWjbW2V5uVtmdpAr33Lhtz2sb9m8WotP6WxvjWxGb4iVsPPM5EGkwq5NCNwb6sn9breK588SiEcBtQEPyK7wXzXBT9QcCZ3S5XWygm#zUC73A7EHiDAxxy29qox4hD5Dyc6fXqStkWjbW2V5uVtmdpAr33Lhtz2sb9m8WotP6WxvjWxGb4iVsPPM5EGkwq5NCNwb6sn9breK588SiEcBtQEPyK7wXzXBT9QcCZ3S5XWygm',
    },
    type: ['VerifiableCredential', 'PermanentResidentCard'],
  },
  vcSubjectData: [
    { name: 'Given Name', value: 'Louis' },
    { name: 'Family Name', value: 'Pasteur' },
    { name: 'Gender', value: 'Male' },
    { name: 'Date of Birth', value: '1958-07-17' },
    { name: 'Country of Birth', value: 'Bahamas' },
    { name: 'Resident Since', value: '2015-01-01' },
  ],
  vpRequest: {
    type: 'QueryByExample',
    credentialQuery: {
      reason: 'Please present your identity document.',
      example: {
        '@context': [
          'https://www.w3.org/2018/credentials/v1',
          'https://w3id.org/citizenship/v1',
          'https://w3id.org/security/bbs/v1',
        ],
        type: ['PermanentResidentCard'],
      },
    },
  },
});
credential.set('VaccinationCertificate', {
  name: 'COVID-19 Vaccination Certificate',
  vc: {
    '@context': [
      'https://www.w3.org/2018/credentials/v1',
      'https://w3id.org/vaccination/v1',
      'https://w3id.org/vc-revocation-list-2020/v1',
      'https://w3id.org/security/bbs/v1',
    ],
    credentialStatus: {
      id: 'https://issuer-vcs.stg.trustbloc.dev/didkey-bbsblssignature2020-bls12381g2/status/1#3',
      revocationListCredential:
        'https://issuer-vcs.stg.trustbloc.dev/didkey-bbsblssignature2020-bls12381g2/status/1',
      revocationListIndex: '3',
      type: 'RevocationList2020Status',
    },
    credentialSubject: {
      administeringCentre: 'FEMA',
      batchNumber: '1183738569',
      countryOfVaccination: 'US',
      dateOfVaccination: '2021-02-01',
      healthProfessional: 'FEMA',
      id: 'did:orb:interim:EiDyiUSqnYPQwMstHKfAY69uU1rbggoUq31edX-LmVD1AQ',
      recipient: {
        birthDate: '1958-07-17',
        familyName: 'Pasteur',
        gender: 'Male',
        givenName: 'Louis',
        type: 'VaccineRecipient',
      },
      type: 'VaccinationEvent',
      vaccine: {
        atcCode: 'J07BX03',
        disease: 'COVID-19',
        marketingAuthorizationHolder: 'Moderna Biotech',
        medicinalProductName: 'COVID-19 Vaccine Moderna',
        type: 'Vaccine',
      },
    },
    description: 'COVID-19 Vaccination Certificate for Mr.Louis Pasteur',
    id: 'http://example.com/9ac3a5da-cd1f-45d1-bfe4-de465c1dd750',
    issuanceDate: '2021-07-06T22:44:27.569613748Z',
    issuer: {
      id: 'did:key:zUC72c7u4BYVmfYinDceXkNAwzPEyuEE23kUmJDjLy8495KH3pjLwFhae1Fww9qxxRdLnS2VNNwni6W3KbYZKsicDtiNNEp76fYWR6HCD8jAz6ihwmLRjcHH6kB294Xfg1SL1qQ',
      name: 'didkey-bbsblssignature2020-bls12381g2',
    },
    name: 'COVID-19 Vaccination Certificate',
    proof: {
      created: '2021-07-06T22:45:54.16001986Z',
      proofPurpose: 'assertionMethod',
      proofValue:
        'sX3b4Y4kpdzHtxGVUkTJ_xKJC1IWCk8EMahJsfFK05vb56QTwBR6cqnlyRXic-AzUwTNH8jKShrj8DP0X0UHX9JvRowQcIwR_99PDvk1nuIgDYh1f3SmKCjRmd5eX6H7E2IQTSiodv9FU1AgTA6otg',
      type: 'BbsBlsSignature2020',
      verificationMethod:
        'did:key:zUC72c7u4BYVmfYinDceXkNAwzPEyuEE23kUmJDjLy8495KH3pjLwFhae1Fww9qxxRdLnS2VNNwni6W3KbYZKsicDtiNNEp76fYWR6HCD8jAz6ihwmLRjcHH6kB294Xfg1SL1qQ#zUC72c7u4BYVmfYinDceXkNAwzPEyuEE23kUmJDjLy8495KH3pjLwFhae1Fww9qxxRdLnS2VNNwni6W3KbYZKsicDtiNNEp76fYWR6HCD8jAz6ihwmLRjcHH6kB294Xfg1SL1qQ',
    },
    type: ['VerifiableCredential', 'VaccinationCertificate'],
  },
  vcSubjectData: [
    { name: 'Given Name', value: 'Louis' },
    { name: 'Family Name', value: 'Pasteur' },
    { name: 'Gender', value: 'Male' },
    { name: 'Date of Birth', value: '1958-07-17' },
    { name: 'Administering Centre', value: 'FEMA' },
    { name: 'Batch Number', value: '1183738569' },
    { name: 'Vaccination Country', value: 'US' },
    { name: 'Date of Vaccination', value: '2021-02-01' },
    { name: 'Health Professional', value: 'FEMA' },
    { name: 'Vaccination Code', value: 'J07BX03' },
    { name: 'Product Name', value: 'Moderna Biotech' },
  ],
  vpRequest: {
    type: 'QueryByFrame',
    credentialQuery: {
      reason: 'Please present your proof of vaccination.',
      frame: {
        '@context': [
          'https://www.w3.org/2018/credentials/v1',
          'https://w3id.org/vaccination/v1',
          'https://w3id.org/vc-revocation-list-2020/v1',
          'https://w3id.org/security/bbs/v1',
        ],
        type: ['VerifiableCredential', 'VaccinationCertificate'],
        '@explicit': true,
        issuer: {},
        issuanceDate: {},
        credentialSubject: {
          '@explicit': true,
          type: 'VaccinationEvent',
          countryOfVaccination: {},
          recipient: {
            '@explicit': true,
            type: 'VaccineRecipient',
            givenName: {},
            familyName: {},
          },
        },
      },
      example: {
        '@context': [
          'https://www.w3.org/2018/credentials/v1',
          'https://w3id.org/vaccination/v1',
          'https://w3id.org/vc-revocation-list-2020/v1',
          'https://w3id.org/security/bbs/v1',
        ],
        type: ['VaccinationCertificate'],
      },
    },
  },
});
credential.set('BookingReference', {
  name: 'Taylor Flights Booking Reference',
  vc: {
    '@context': [
      'https://www.w3.org/2018/credentials/v1',
      'https://trustbloc.github.io/context/vc/examples/booking-ref-v1.jsonld',
    ],
    credentialSubject: {
      id: 'did:key:z6MkjRagNiMu91DduvCvgEsqLZDVzrJzFrwahc4tXLt9DoHd',
      issuedBy: 'Taylor Chartered Flights',
      referenceNumber: 'W7X 9T3',
    },
    description: 'Booking reference of Mr.John Smith',
    id: 'http://example.gov/credentials/3732',
    issuanceDate: '2020-03-16T22:37:26.544Z',
    issuer: 'did:key:z6MkjRagNiMu91DduvCvgEsqLZDVzrJzFrwahc4tXLt9DoHd',
    name: 'Taylor Flights Booking Reference',
    proof: {
      challenge: '69b25d39-e87c-4627-ab1a-144c632ca211',
      created: '2021-03-11T10:02:03.519525-05:00',
      domain: 'issuer.service.com',
      jws: 'eyJhbGciOiJFZERTQSIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..Nf89Zh5tzKJiEKSFH0jh4P4NQuZI1S0IpN1HRsSWyvg5N3Cmm2QT1Zoz3I3WGV7pEhLqFbfr8AaOc_RmJWYuDQ',
      proofPurpose: 'authentication',
      type: 'Ed25519Signature2018',
      verificationMethod: 'did:example:xyz#key-1',
    },
    type: ['VerifiableCredential', 'BookingReferenceCredential'],
  },
  vcSubjectData: [
    { name: 'Reference Number', value: 'W7X 9T3' },
    { name: 'Issued By', value: 'Taylor Chartered Flights' },
  ],
  vpRequest: {
    type: 'QueryByExample',
    credentialQuery: {
      reason: 'Please present your booking reference.',
      example: {
        '@context': [
          'https://www.w3.org/2018/credentials/v1',
          'https://trustbloc.github.io/context/vc/examples/booking-ref-v1.jsonld',
        ],
        type: ['BookingReferenceCredential'],
      },
    },
  },
});
credential.set('GovernanceCredential', {
  name: 'Governance Credential',
  vc: {
    '@context': [
      'https://www.w3.org/2018/credentials/v1',
      'https://trustbloc.github.io/context/governance/context.jsonld',
      'https://w3id.org/vc-revocation-list-2020/v1',
    ],
    credentialSubject: {
      description: 'Sample governance framework for the TrustBloc sandbox.',
      geos: 'Canadian',
      jurisdictions: 'ca',
      roles: 'accreditor',
      topics: 'banking',
      version: '1.0',
    },
    id: 'https://example.com/governance',
    issuanceDate: '2021-07-06T17:20:57.736460722Z',
    issuer: 'did:orb:interim:EiAbjeimRkgIIswIPpHpWO3JeXZg3eeuVjlzzCeAjgOC-g',
    name: 'Governance Credential',
    proof: {
      created: '2021-07-06T17:20:57.851590383Z',
      jws: 'eyJhbGciOiJFZERTQSIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..92EdftDqyPKq4HYjjwRrtYPUBVWukVxbhweNvque-QfMSkL7FKww-5Q25hq5uowa-6mFiV3UvB-9xC2x5s-eCg',
      proofPurpose: 'assertionMethod',
      type: 'Ed25519Signature2018',
      verificationMethod:
        'did:orb:interim:EiAbjeimRkgIIswIPpHpWO3JeXZg3eeuVjlzzCeAjgOC-g#a5hRq_BgA-TU0-RUsuhoIo3qQhrTEt2iKO9MPDeZ2fc',
    },
    type: ['VerifiableCredential', 'GovernanceCredential'],
  },
  vcSubjectData: [
    { name: 'Description', value: 'Sample governance framework' },
    { name: 'Geos', value: 'Canadian' },
    { name: 'Jurisdictions', value: 'ca' },
    { name: 'Roles', value: 'accreditor' },
    { name: 'Topics', value: 'banking' },
    { name: 'Version', value: '1.0' },
  ],
  vpRequest: {
    type: 'QueryByExample',
    credentialQuery: {
      reason: 'Please present your Governance Credential.',
      example: {
        '@context': [
          'https://www.w3.org/2018/credentials/v1',
          'https://trustbloc.github.io/context/governance/context.jsonld',
        ],
        type: ['GovernanceCredential'],
      },
    },
  },
});

describe('TrustBloc Wallet - Store/Share credential flow (CHAPI)', function () {
  let ctx;

  // runs once before the first test in this block
  before(async function () {
    await browser.reloadSession();
    await browser.maximizeWindow();
    ctx = {
      email: `ui-aut-chapi-${new Date().getTime()}@test.com`,
    };
  });

  it(`User Sign up`, async function () {
    // 1. Navigate to Wallet Website
    await browser.navigateTo(browser.config.walletURL);

    // 2. Initialize Wallet (register/sign-up/etc.)
    await wallet.signUp(ctx, browser.config.isCHAPIEnabled);

    await wallet.waitForCredentials();
  });

  it(`Create Orb DID`, async function () {
    await wallet.createOrbDID();
  });

  it(`Import DID Key (JWK key format)`, async function () {
    await wallet.importDID({ method: 'key' });
  });

  // TODO add case to import DID with Base58 key format

  it(`Update Digital Identity preferences`, async function () {
    await wallet.updatePreferences();
  });

  it(`User performs DID Auth with mock issuer`, async function () {
    // mock issuer (wallet page with sample requests)
    await browser.navigateTo(browser.config.chapiDemoURL);

    const didAuthBtn = await $('#didauth');
    await didAuthBtn.waitForClickable();
    await didAuthBtn.click();

    const storeButton = await $('#get-btn');
    await storeButton.waitForClickable();
    await storeButton.click();

    await chapi.chooseWallet({
      name: browser.config.walletName,
    });

    await wallet.authenticate(ctx);
    await browser.switchToFrame(null);

    const didAuthSuccessMsg = await $('div*=Successfully got response from wallet.');
    await didAuthSuccessMsg.waitForExist();
  });

  it(`User stores credential from mock issuer`, async function () {
    for (const [key, value] of credential.entries()) {
      // mock issuer (wallet page with sample requests)
      await browser.navigateTo(browser.config.chapiDemoURL);

      console.log('save vc : start ', key);

      const vcSampleBtn = await $('#store-vc-sample-1');
      await vcSampleBtn.waitForClickable();
      await vcSampleBtn.click();

      const vprs = {
        '@context': ['https://www.w3.org/2018/credentials/v1'],
        type: 'VerifiablePresentation',
        verifiableCredential: [value.vc],
      };

      const sampleText = await $('#vcDataTextArea');
      await sampleText.waitForExist();
      await sampleText.setValue(vprs);

      const storeButton = await $('#store-btn');
      await storeButton.waitForClickable();
      await storeButton.click();

      await chapi.chooseWallet({
        name: browser.config.walletName,
      });

      await wallet.validateCredentialDetails(value.vcSubjectData);

      await wallet.storeCredentials(ctx);
      await browser.switchToFrame(null);

      const storeSuccessMsg = await $(
        'div*=Successfully stored verifiable presentation to wallet.'
      );
      await expect(storeSuccessMsg).toExist();

      console.log('save vc : end ', key);
    }
  });

  it(`User validates the saved credential from mock issuer`, async function () {
    await browser.navigateTo(`${browser.config.walletURL}/credentials`);

    for (const [key, value] of credential.entries()) {
      console.log('validate vc in wallet : start ', key);

      const vcName = await $('span*=' + value.name);
      await vcName.waitForClickable();
      await vcName.click();

      await wallet.validateCredentialDetails(value.vcSubjectData);

      const credTab = await $('span*=Credentials');
      await credTab.waitForClickable();
      await credTab.click();

      console.log('validate vc in wallet : end ', key);
    }
  });

  it(`User Sign Out`, async function () {
    await wallet.signOut(ctx);
  });

  it(`User Sign in`, async function () {
    // Sign In to the registered Wallet
    await wallet.signIn(ctx.email);

    await wallet.waitForCredentials();
  });

  it(`User validates the saved credential from mock issuer (after sign-in)`, async function () {
    for (const [key, value] of credential.entries()) {
      console.log('validate vc in wallet : start ', key);

      const vcName = await $('span*=' + value.name);
      await vcName.waitForClickable();
      await vcName.click();

      await wallet.validateCredentialDetails(value.vcSubjectData);

      const credTab = await $('span*=Credentials');
      await credTab.waitForClickable();
      await credTab.click();

      console.log('validate vc in wallet : end ', key);
    }
  });

  it(`User shares the saved credential with mock verifier`, async function () {
    // mock verifier (wallet page with sample requests)
    await browser.navigateTo(browser.config.chapiDemoURL);

    for (const [key, value] of credential.entries()) {
      console.log('share vc : start ', key);

      const vcSampleBtn = await $('#req-vp');
      await vcSampleBtn.waitForClickable();
      await vcSampleBtn.click();

      const vprs = {
        web: {
          VerifiablePresentation: {
            query: [value.vpRequest],
            challenge: '4ada78c9-d58f-4e93-a039-cfa6c8999b97',
            domain: 'example.com',
          },
        },
      };

      const sampleText = await $('#vcDataTextArea');
      await sampleText.clearValue();
      await sampleText.addValue(vprs);

      const getButton = await $('#get-btn');
      await getButton.waitForClickable();
      await getButton.click();

      await chapi.chooseWallet({
        name: browser.config.walletName,
      });

      // TODO https://github.com/trustbloc/wallet/issues/1124 VC Name mismatch between credentials screen and CHAPI share
      // const vcName = await $('span*=' + value.name);
      // await vcName.waitForClickable();
      // await vcName.click();

      // await wallet.validateCredentialDetails(value.vcSubjectData);

      await wallet.presentCredentials(ctx);
      await browser.switchToFrame(null);

      const getSuccessMsg = await $('div*=Successfully got response from wallet.');
      await expect(getSuccessMsg).toExist();

      console.log('share vc : end ', key);
    }
  });
  it(`User deletes the saved credential`, async function () {
    await browser.navigateTo(`${browser.config.walletURL}/credentials`);

    let numOfCreds = credential.size;
    const credentialsList = await $('#credentials-list');
    await credentialsList.waitForExist();
    await expect(credentialsList).toHaveChildren(numOfCreds);

    for (const [key, value] of credential.entries()) {
      console.log('delete vc in wallet : start ', key);

      const vcName = await $('span*=' + value.name);
      await vcName.waitForClickable();
      await vcName.click();

      await wallet.validateCredentialDetails(value.vcSubjectData);
      console.log('delete vc in wallet : end ', key);
      await wallet.deleteCredential(numOfCreds);

      --numOfCreds;
      if (numOfCreds > 0) await expect(credentialsList).toHaveChildren(numOfCreds);
    }

    const noCredsContainer = await $('#loaded-credentials-container');
    await expect(noCredsContainer).toExist();
  });

  it(`User Sign Out`, async function () {
    await wallet.signOut(ctx);
  });

  it(`User changes locale`, async function () {
    await wallet.changeLocale();
  });
});
