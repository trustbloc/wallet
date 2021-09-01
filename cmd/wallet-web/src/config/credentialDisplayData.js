/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
export default {
  PermanentResidentCard: {
    properties: {
      birthCountry: {
        path: ['$.credentialSubject.birthCountry'],
        label: 'Country of Birth',
      },
      birthDate: {
        path: ['$.credentialSubject.birthDate'],
        label: 'Date of birth',
        type: 'date',
        format: 'yyyy-mm-dd',
      },
      familyName: {
        path: ['$.credentialSubject.familyName'],
        label: 'Family Name',
      },
      givenName: {
        path: ['$.credentialSubject.givenName'],
        label: 'Given Name',
      },
      residentSince: {
        path: ['$.credentialSubject.residentSince'],
        label: 'Resident Since',
        type: 'date',
        format: 'yyyy-mm-dd',
      },
      image: {
        path: ['$.credentialSubject.image'],
        label: 'Photo',
        type: 'image',
        format: '',
      },
    },
    icon: 'credential--generic-icon.svg',
    title: {
      path: ['$.name'],
      fallback: 'Permanent Resident Card',
    },
    subTitle: {
      path: ['$.description'],
      fallback: 'Permanent Resident Card',
    },
  },
  UniversityDegreeCredential: {
    properties: {
      type: {
        path: ['$.credentialSubject.degree.type'],
        label: 'Degree Category',
      },
      name: {
        path: ['$.credentialSubject.degree.name'],
        label: 'Degree Name',
      },
    },
    icon: 'credential--school-icon.svg',
    title: {
      path: ['$.name'],
      fallback: 'Bachelor Degree',
    },
  },
  VaccinationCertificate: {
    properties: {
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
      RecipientBirthDate: {
        path: ['$.credentialSubject.recipient.birthDate'],
        label: 'Date of Birth',
        type: 'date',
        format: 'yyyy-mm-dd',
      },
      RecipientFamilyName: {
        path: ['$.credentialSubject.recipient.familyName'],
        label: 'Family Name',
      },
      RecipientGender: {
        path: ['$.credentialSubject.recipient.gender'],
        label: 'Gender',
      },
      RecipientGivenName: {
        path: ['$.credentialSubject.recipient.givenName'],
        label: 'Given Name',
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
    icon: 'credential--generic-icon.svg',
    title: {
      path: ['$.name'],
      fallback: 'Vaccination Certificate',
    },
  },
  TravelCard: {
    properties: {
      name: {
        path: ['$.credentialSubject.name'],
        label: 'Cardholder Name',
      },
      sex: {
        path: ['$.credentialSubject.sex'],
        label: 'Gender',
      },
      country: {
        path: ['$.credentialSubject.country'],
        label: 'Country',
      },
      dob: {
        path: ['$.credentialSubject.dob'],
        label: 'Date of Birth',
        type: 'date',
        format: 'yyyy-mm-dd',
      },
      travelcardid: {
        path: ['$.credentialSubject.travelcardid'],
        label: 'Travel Card No.',
      },
      issuedate: {
        path: ['$.credentialSubject.issuedate'],
        label: 'Date of Issue',
        type: 'date',
        format: 'yyyy-mm-dd',
      },
      cardexpires: {
        path: ['$.credentialSubject.cardexpires'],
        label: 'Date of Expiry',
        type: 'date',
        format: 'yyyy-mm-dd',
      },
    },
    icon: 'credential--generic-icon.svg',
    title: {
      path: ['$.name'],
      fallback: 'Travel Card',
    },
  },
  StudentCard: {
    properties: {
      name: {
        path: ['$.credentialSubject.name'],
        label: 'Student Name',
      },
      email: {
        path: ['$.credentialSubject.email'],
        label: 'Email',
      },
      studentid: {
        path: ['$.credentialSubject.studentid'],
        label: 'Student No.',
      },
      university: {
        path: ['$.credentialSubject.university'],
        label: 'University',
      },
      semester: {
        path: ['$.credentialSubject.semester'],
        label: 'Semester',
      },
    },
    icon: 'credential--generic-icon.svg',
    title: {
      path: ['$.name'],
      fallback: 'Student Card',
    },
  },
  mdL: {
    properties: {
      familyName: {
        path: ['$.credentialSubject.family_name'],
        label: 'Family Name',
      },
      givenName: {
        path: ['$.credentialSubject.given_name'],
        label: 'Given Name',
      },
      birthdate: {
        path: ['$.credentialSubject.birthdate'],
        label: 'Date of Birth',
        type: 'date',
        format: 'yyyy-mm-dd',
      },
      driving_privileges: {
        path: ['$.credentialSubject.driving_privileges'],
        label: 'Class',
      },
      issuing_authority: {
        path: ['$.credentialSubject.issuing_authority'],
        label: 'Issued by',
      },
      resident_address: {
        path: ['$.credentialSubject.resident_address'],
        label: 'Address',
      },
      resident_city: {
        path: ['$.credentialSubject.resident_city'],
        label: 'City',
      },
      resident_postal_code: {
        path: ['$.credentialSubject.resident_postal_code'],
        label: 'Postal Code',
      },
      issue_date: {
        path: ['$.credentialSubject.issue_date'],
        label: 'Date of Issue',
        type: 'date',
        format: 'yyyy-mm-dd',
      },
    },
    icon: 'credential--generic-icon.svg',
    title: {
      path: ['$.name'],
      fallback: 'Driving License',
    },
  },
  fallback: {
    properties: {},
    icon: 'credential--generic-icon.svg',
    title: {
      path: ['$.name'],
      fallback: 'Verifiable Credential',
    },
  },
};
