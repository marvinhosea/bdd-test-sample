Feature: Search packages on pkg.go.dev
  Navigate go.dev and change path to pkg.go.dev to search pkg.go.dev to search for packages

  Background: Load Go.dev
    Given I visit "https://go.dev"

  Scenario Outline: Search package on https://pkg.go.dev
    When I navigate to https://pkg.go.dev by clicking packages on menu
    And I enter "<package>" package name in the search
    And I press search button
    Then I should see a search page with "<packageUrl>" package

    Examples:
      | package     | packageUrl                    |
      | godog       | github.com/cucumber/godog     |
      | dotenv      | github.com/joho/godotenv      |
      | marvinhosea | github.com/marvinhosea/filter |