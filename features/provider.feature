Feature: UpCloud DevPod Provider
  As a DevPod user
  I want to use UpCloud as my cloud provider
  So that I can create and manage development environments on UpCloud infrastructure

  Background:
    Given I have valid UpCloud API credentials
    And the provider is configured with required options

  Scenario: Initialize provider with valid credentials
    When I run the init command
    Then the provider should validate the credentials
    And return a success status

  Scenario: Create a new UpCloud server
    When I run the create command
    Then a new server should be created in UpCloud
    And the server should be accessible via SSH
    And the status should return "Running"

  Scenario: Stop a running server
    Given I have a running UpCloud server
    When I run the stop command
    Then the server should be stopped
    And the status should return "Stopped"

  Scenario: Start a stopped server
    Given I have a stopped UpCloud server
    When I run the start command
    Then the server should be started
    And the status should return "Running"

  Scenario: Delete a server
    Given I have an existing UpCloud server
    When I run the delete command
    Then the server should be removed from UpCloud
    And the status should return "NotFound"

  Scenario: Execute command on server
    Given I have a running UpCloud server
    When I execute a command on the server
    Then the command should run successfully
    And I should see the command output