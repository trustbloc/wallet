/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

export default {
  PermanentResidentCard: {
    id: ['$.id'],
    brandColor: 'boatBlue',
    issuanceDate: ['$.issuanceDate'],
    properties: {
      image: {
        path: ['$.credentialSubject.image'],
        label: 'Photo',
        type: 'image',
        format: '',
      },
      givenName: {
        path: ['$.credentialSubject.givenName'],
        label: 'Given Name',
      },
      familyName: {
        path: ['$.credentialSubject.familyName'],
        label: 'Family Name',
      },
      gender: {
        path: ['$.credentialSubject.gender'],
        label: 'Gender',
      },
      birthDate: {
        path: ['$.credentialSubject.birthDate'],
        label: 'Date of birth',
        type: 'date',
        format: 'yyyy-mm-dd',
      },
      birthCountry: {
        path: ['$.credentialSubject.birthCountry'],
        label: 'Country of Birth',
      },
      residentSince: {
        path: ['$.credentialSubject.residentSince'],
        label: 'Resident Since',
        type: 'date',
        format: 'yyyy-mm-dd',
      },
    },
    icon: 'credential--uscis-icon.svg',
    title: {
      path: ['$.name'],
      fallback: 'Permanent Resident Card',
    },
    subTitle: {
      path: ['$.description'],
      fallback: 'Permanent Resident Card',
    },
  },
  VaccinationCertificate: {
    id: ['$.id'],
    brandColor: 'green',
    issuanceDate: ['$.issuanceDate'],
    properties: {
      RecipientFamilyName: {
        path: ['$.credentialSubject..familyName'],
        label: 'Family Name',
      },
      RecipientGivenName: {
        path: ['$.credentialSubject..givenName'],
        label: 'Given Name',
      },
      RecipientGender: {
        path: ['$.credentialSubject..gender'],
        label: 'Gender',
      },
      RecipientBirthDate: {
        path: ['$.credentialSubject..birthDate'],
        label: 'Date of Birth',
        type: 'date',
        format: 'yyyy-mm-dd',
      },
      administeringCentre: {
        path: ['$.credentialSubject.administeringCentre'],
        label: 'Administering Centre',
      },
      batchNumber: {
        path: ['$.credentialSubject.batchNumber'],
        label: 'Batch Number',
      },
      countryOfVaccination: {
        path: ['$.credentialSubject.countryOfVaccination'],
        label: 'Vaccination Country',
      },
      dateOfVaccination: {
        path: ['$.credentialSubject.dateOfVaccination'],
        label: 'Date of Vaccination',
      },
      healthProfessional: {
        path: ['$.credentialSubject.healthProfessional'],
        label: 'Health Professional',
      },
      VaccineCode: {
        path: ['$.credentialSubject.vaccine.atcCode'],
        label: 'Vaccination Code',
      },
      VaccineProductName: {
        path: ['$.credentialSubject.vaccine.medicinalProductName'],
        label: 'Product Name',
      },
      VaccineAuthorizationHolder: {
        path: ['$.credentialSubject.vaccine.marketingAuthorizationHolder'],
        label: 'Product Name',
      },
    },
    icon: 'credential--vaccination-icon.svg',
    title: {
      path: ['$.name'],
      fallback: 'Vaccination Certificate',
    },
  },
  BookingReferenceCredential: {
    id: ['$.id'],
    brandColor: 'cobalt',
    issuanceDate: ['$.issuanceDate'],
    properties: {
      issuedBy: {
        path: ['$.credentialSubject.issuedBy'],
        label: 'Issued By',
      },
      referenceNumber: {
        path: ['$.credentialSubject.referenceNumber'],
        label: 'Reference Number',
      },
    },
    icon: 'credential--flight-icon.svg',
    title: {
      path: ['$.name'],
      fallback: 'Booking Reference',
    },
  },
  fallback: {
    id: ['$.id'],
    issuanceDate: ['$.issuanceDate'],
    properties: {},
    icon: 'credential--generic-icon.svg',
    title: {
      path: ['$.name'],
      fallback: 'Verifiable Credential',
    },
  },
};
