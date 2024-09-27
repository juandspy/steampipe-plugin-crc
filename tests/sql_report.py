"""
This script extracts SQL commands from a README file, runs them using Steampipe, and generates a report.
"""


import re
import subprocess
from datetime import datetime

def extract_sql_commands(file_path):
    with open(file_path, 'r') as file:
        content = file.read()
    
    # Find all SQL code blocks
    sql_blocks = re.findall(r'```sql\n(.*?)```', content, re.DOTALL)
    
    # Clean up the SQL commands
    sql_commands = [cmd.strip() for block in sql_blocks for cmd in block.split(';') if cmd.strip()]
    return sql_commands

def run_sql_command(command):
    print("-" * 30)
    print(f"Running command: {command}")
    try:
        result = subprocess.run(['steampipe', 'query', command], 
                                capture_output=True, text=True, timeout=300)
        return result.returncode == 0, result.stdout, result.stderr
    except subprocess.TimeoutExpired:
        return False, "", "Command timed out after 5 minutes"

def generate_report(results):
    success_count = sum(1 for r in results if r['success'])
    failure_count = len(results) - success_count
    
    report = f"SQL Command Test Report\n"
    report += f"Generated on: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n\n"
    report += f"Total commands tested: {len(results)}\n"
    report += f"Successful: {success_count}\n"
    report += f"Failed: {failure_count}\n\n"
    
    for i, result in enumerate(results, 1):
        report += f"Command {i}:\n"
        report += f"SQL: {result['command']}\n"
        report += f"Success: {'Yes' if result['success'] else 'No'}\n"
        if result['success']:
            report += f"Output: {result['output']}\n"
        else:
            report += f"Error: {result['error']}\n"

        report += "-" * 30
        report += "\n"
    
    return report

def main():
    readme_path = 'README.md'
    sql_commands = extract_sql_commands(readme_path)
    
    results = []
    for command in sql_commands:
        success, output, error = run_sql_command(command)
        results.append({
            'command': command,
            'success': success,
            'output': output,
            'error': error
        })
    
    report = generate_report(results)
    
    with open('sql_test_report.txt', 'w') as f:
        f.write(report)
    
    print("Test complete. Report generated in sql_test_report.txt")

if __name__ == "__main__":
    main()