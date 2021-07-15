#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

@all
@signup
Feature: Signup

#issue-838 will fix bdd test as all scenarios are tied to each other
#Scenario: User Signup
#    When the user clicks on the Signup button
#     And the user is redirected to the OIDC provider
#    And the user is authenticated
#    And the user consents to sharing their identity data
#    Then the user is redirected to the wallet's dashboard
#    And the user can retrieve their profile
#
# Scenario: User logout
#    When the user is logged in
#    And the user logs out
#    Then the user cannot access their profile
#
#  Scenario: Device registration and login through device
#    When the user is logged in
#    And the user registers a device
#    And the user logs out
#    Then the user logs in through their device
#    And the user can retrieve their profile
