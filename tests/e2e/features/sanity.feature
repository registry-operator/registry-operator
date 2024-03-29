Feature: Sanity
  Set of scenarios for checking if the environment is set-up correctly

  Scenario: Testing fixture
    Given kubernetes cluster
    Then e2e namespace should exist
