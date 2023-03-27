import sys
from lxml import etree

def parse_outcomes(file_path):
    outcomes = []

    with open(file_path, 'r') as file:
        xml_string = file.read()

    root = etree.fromstring(xml_string)
    for outcome_elem in root.findall('outcome'):
        outcome = {
            'home_team': outcome_elem.find('homeTeam').text,
            'out_team': outcome_elem.find('outTeam').text,
            'home_score': int(outcome_elem.find('homeScore').text),
            'out_score': int(outcome_elem.find('outScore').text),
            'date': outcome_elem.find('date').text,
        }
        outcomes.append(outcome)
    return outcomes

def generate_teams_insert_sql(outcomes):
    teams = set()

    for outcome in outcomes:
        teams.add(outcome['home_team'])
        teams.add(outcome['out_team'])

    sql_statements = []

    for team in teams:
        statement = f"INSERT INTO teams (name) VALUES ('{team}');"
        sql_statements.append(statement)

    return sql_statements

def generate_matches_insert_sql(outcomes):
    sql_statements = []

    for outcome in outcomes:
        statement = f"INSERT INTO matches (home_team, away_team, home_goals, away_goals, date) VALUES ((SELECT id FROM teams WHERE name = '{outcome['home_team']}'), (SELECT id FROM teams WHERE name = '{outcome['out_team']}'), {outcome['home_score']}, {outcome['out_score']}, '{outcome['date']}');"
        sql_statements.append(statement)

    return sql_statements

if __name__ == '__main__':
    if len(sys.argv) != 2:
        print("Usage: python script.py <path_to_xml_file>")
        sys.exit(1)

    file_path = sys.argv[1]
    outcomes = parse_outcomes(file_path)

    teams_insert_sql = generate_teams_insert_sql(outcomes)
    matches_insert_sql = generate_matches_insert_sql(outcomes)

    print("-- Populate teams table")
    for statement in teams_insert_sql:
        print(statement)

    print("\n-- Populate matches table")
    for statement in matches_insert_sql:
        print(statement)
