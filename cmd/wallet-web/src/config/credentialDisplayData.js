/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
export default {
  PermanentResidentCard: {
    id: ['$.id'],
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
    id: ['$.id'],
    issuanceDate: ['$.issuanceDate'],
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
    id: ['$.id'],
    issuanceDate: ['$.issuanceDate'],
    properties: {
      RecipientFamilyName: {
        path: ['$.credentialSubject.recipient.familyName'],
        label: 'Family Name',
      },
      RecipientGivenName: {
        path: ['$.credentialSubject.recipient.givenName'],
        label: 'Given Name',
      },
      RecipientGender: {
        path: ['$.credentialSubject.recipient.gender'],
        label: 'Gender',
      },
      RecipientBirthDate: {
        path: ['$.credentialSubject.recipient.birthDate'],
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
    icon: 'credential--generic-icon.svg',
    title: {
      path: ['$.name'],
      fallback: 'Vaccination Certificate',
    },
  },
  TravelCard: {
    id: ['$.id'],
    issuanceDate: ['$.issuanceDate'],
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
    id: ['$.id'],
    issuanceDate: ['$.issuanceDate'],
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
    id: ['$.id'],
    issuanceDate: ['$.issuanceDate'],
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
  CrudeProductCredential: {
    id: ['$.id'],
    issuanceDate: ['$.issuanceDate'],
    properties: {
      name: {
        path: ['$.credentialSubject.name'],
        label: 'Name',
      },
      description: {
        path: ['$.credentialSubject.description'],
        label: 'Description',
      },
      productionDate: {
        path: ['$.credentialSubject.productionDate'],
        label: 'Date of Production',
        type: 'date',
        format: 'yyyy-mm-dd',
      },
      category: {
        path: ['$.credentialSubject.category'],
        label: 'Category',
      },
      hsCode: {
        path: ['$.credentialSubject.hsCode'],
        label: 'Code',
      },
      volume: {
        path: ['$.credentialSubject.volume'],
        label: 'Volume',
      },
      minimumQuantity: {
        path: ['$.credentialSubject.physicalSpecs.minimumQuantity'],
        label: 'Min Quantity',
      },
      uom: {
        path: ['$.credentialSubject.physicalSpecs.uom'],
        label: 'Unit of Measurement',
      },
      pourPoint: {
        path: ['$.credentialSubject.physicalSpecs.pourPoint'],
        label: 'Pour Point',
      },
      vapourPressure: {
        path: ['$.credentialSubject.physicalSpecs.vapourPressure'],
        label: 'Vapour Pressure',
      },
      density: {
        path: ['$.credentialSubject.physicalSpecs.density'],
        label: 'Density',
      },
      microCarbonResidue: {
        path: ['$.credentialSubject.chemicalSpecs.microCarbonResidue'],
        label: 'Carbon Residue',
      },
    },
    icon: 'credential--generic-icon.svg',
    title: {
      path: ['$.name'],
      fallback: 'Crude Product Credential',
    },
  },
  CertifiedMillTestReport: {
    id: ['$.id'],
    issuanceDate: ['$.issuanceDate'],
    properties: {
      companyBrandMark: {
        path: ['$.credentialSubject.cmtr.companyBrandMark'],
        label: 'BrandMark',
        type: 'image',
      },
      companyName: {
        path: ['$.credentialSubject.cmtr.companyName'],
        label: 'Name',
      },
      companyAddress: {
        path: ['$.credentialSubject.cmtr.companyAddress'],
        label: 'Address',
      },
      companyPhoneNumber: {
        path: ['$.credentialSubject.cmtr.companyPhoneNumber'],
        label: 'Phone Number',
      },
      companyContactPersonName: {
        path: ['$.credentialSubject.cmtr.companyContactPersonName'],
        label: 'Contact Person',
      },
      invoiceNumber: {
        path: ['$.credentialSubject.cmtr.invoiceNumber'],
        label: 'Invoice Number',
      },
      authorizingPartyName: {
        path: ['$.credentialSubject.cmtr.authorizingPartyName'],
        label: 'Authorizing Party Name',
      },
      authorizingPartyDate: {
        path: ['$.credentialSubject.cmtr.authorizingPartyDate'],
        label: 'Authorizing Party Date',
        type: 'date',
        format: 'yyyy-mm-dd',
      },
      manufacturerLocationCompanyName: {
        path: ['$.credentialSubject.cmtr.manufacturerLocationCompanyName'],
        label: 'Manufacturer Name',
      },
      manufacturerLocationStreetAddress: {
        path: ['$.credentialSubject.cmtr.manufacturerLocationStreetAddress'],
        label: 'Manufacturer Street Address',
      },
      manufacturerLocationAddressLocality: {
        path: ['$.credentialSubject.cmtr.manufacturerLocationAddressLocality'],
        label: 'Manufacturer City',
      },
      manufacturerLocationAddressCountry: {
        path: ['$.credentialSubject.cmtr.manufacturerLocationAddressCountry'],
        label: 'Manufacturer Location',
      },
      customerLocationCompanyName: {
        path: ['$.credentialSubject.cmtr.customerLocationCompanyName'],
        label: 'Customer Location',
      },
      customerLocationStreetAddress: {
        path: ['$.credentialSubject.cmtr.customerLocationStreetAddress'],
        label: 'Customer Address',
      },
      customerLocationAddressLocality: {
        path: ['$.credentialSubject.cmtr.customerLocationAddressLocality'],
        label: 'Customer City',
      },
      customerLocationAddressCountry: {
        path: ['$.credentialSubject.cmtr.customerLocationAddressCountry'],
        label: 'Customer Country',
      },
    },
    icon: 'credential--generic-icon.svg',
    title: {
      path: ['$.name'],
      fallback: 'Certified Mill Test Report',
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
