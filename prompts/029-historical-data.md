# Historical data

> Awesome!
>
> I have an xml file with historical data. Can you create a Python script for me that reads this xml and generates SQL that I can run against the database? Here's what the xml looks like:
>
> <outcomes><outcome><homeTeam>Jong AZ</homeTeam><outTeam>TOP Oss</outTeam><homeScore>1</homeScore><outScore>0</outScore><date>2021-08-23</date></outcome><outcome><homeTeam>Jong PSV</homeTeam><outTeam>FC Eindhoven</outTeam><homeScore>3</homeScore><outScore>2</outScore><date>2021-08-23</date></outcome><outcome><homeTeam>VVV-Venlo</homeTeam><outTeam>Almere City FC</outTeam><homeScore>2</homeScore><outScore>1</outScore><date>2021-08-21</date></outcome><outcome><homeTeam>ADO Den Haag</homeTeam><outTeam>FC Emmen</outTeam><homeScore>1</homeScore><outScore>2</outScore><date>2021-08-20</date></outcome>

Here's a Python script that reads the provided XML and generates SQL statements to insert the data into your database. Please make sure to install the required libraries with pip install lxml before running the script.

```python
import sys
from lxml import etree

def parse_outcomes(xml_string):
    outcomes = []

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

def generate_sql(outcomes):
    sql_statements = []

    for outcome in outcomes:
        statement = f"INSERT INTO matches (home_team_id, away_team_id, home_goals, away_goals, match_date) VALUES ((SELECT id FROM teams WHERE name = '{outcome['home_team']}'), (SELECT id FROM teams WHERE name = '{outcome['out_team']}'), {outcome['home_score']}, {outcome['out_score']}, '{outcome['date']}');"
        sql_statements.append(statement)

    return sql_statements

if __name__ == '__main__':
    xml_string = '''
    <outcomes>
        <!-- Your XML data here -->
    </outcomes>
    '''

    outcomes = parse_outcomes(xml_string)
    sql_statements = generate_sql(outcomes)

    for statement in sql_statements:
        print(statement)
```

Replace the <!-- Your XML data here --> with the XML data you provided. When you run this script, it will print the generated SQL statements that you can use to insert the data into your database.

Please note that this script assumes you have a matches table in your database, with columns home_team_id, away_team_id, home_goals, away_goals, and match_date. The teams table should have id and name columns. Adjust the table and column names as needed to match your database schema.

## Follow-up prompt

> Let's update the script in 2 ways:
>
> 1. It should read the xml from a file, because it's too big to paste into the script.
>
> 2. Let's assume the teams table is empty, so that needs to be populated as well.

Here's an updated Python script that reads the XML from a file and also generates SQL statements to populate the teams table. Save the script in a .py file and run it with the XML file path as a command-line argument.

```python
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
        statement = f"INSERT INTO matches (home_team_id, away_team_id, home_goals, away_goals, match_date) VALUES ((SELECT id FROM teams WHERE name = '{outcome['home_team']}'), (SELECT id FROM teams WHERE name = '{outcome['out_team']}'), {outcome['home_score']}, {outcome['out_score']}, '{outcome['date']}');"
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
```

To run the script, provide the XML file path as a command-line argument:

```bash
python script.py path/to/your/xml_file.xml
```

The script will print the generated SQL statements for populating the teams and matches tables. Make sure the table and column names in the script match your database schema.

## Note

ChatGPT got the column names for the matches table wrong, but that was easy enough to fix.
