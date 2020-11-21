/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

\! echo "Configuring MySQL users...";

/*
Hydra
*/
CREATE USER 'hydra'@'%' IDENTIFIED BY 'hydra-pwd';
CREATE DATABASE hydra;
GRANT ALL PRIVILEGES ON hydra.* TO 'hydra'@'%';

/*
Demo Hydra
*/
CREATE USER 'demohydra'@'%' IDENTIFIED BY 'demohydra-pwd';
CREATE DATABASE demohydra;
GRANT ALL PRIVILEGES ON demohydra.* TO 'demohydra'@'%';

/*
demo auth rest
*/
CREATE USER 'authrest'@'%' IDENTIFIED BY 'authrest-secret-pw';
GRANT ALL PRIVILEGES ON `authrest\_%` . * TO 'authrest'@'%';

/*
demo auth rest hydra
*/
CREATE USER 'authresthydra'@'%' IDENTIFIED BY 'authresthydra-secret-pw';
CREATE DATABASE authresthydra;
GRANT ALL PRIVILEGES ON authresthydra.* TO 'authresthydra'@'%';

/*
BDD Hydra
*/
CREATE USER 'bddhydra'@'%' IDENTIFIED BY 'bddhydra-pwd';
CREATE DATABASE bddhydra;
GRANT ALL PRIVILEGES ON bddhydra.* TO 'bddhydra'@'%';

/*
bdd auth rest
*/
CREATE USER 'bddauthrest'@'%' IDENTIFIED BY 'bddauthrest-secret-pw';
GRANT ALL PRIVILEGES ON `authrest\_%` . * TO 'bddauthrest'@'%';

/*
bd auth rest hydra
*/
CREATE USER 'bddauthresthydra'@'%' IDENTIFIED BY 'bddauthresthydra-secret-pw';
CREATE DATABASE bddauthresthydra;
GRANT ALL PRIVILEGES ON bddauthresthydra.* TO 'bddauthresthydra'@'%';
