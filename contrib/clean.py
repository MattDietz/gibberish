import csv
import re


def clean_data(path, output):
    r = re.compile(r'=+[0-9]+')
    out = open(output, 'w')
    try:
        with open(path, 'r') as f:
            reader = csv.reader(f)
            for i, row in enumerate(reader):
                if i > 195:
                    if len(row) < 5:
                        continue

                    line = r.sub('', row[4])
                    if line.startswith("'"):
                        line = line[1:]
                    if line.endswith("'"):
                        line = line[:-1]
                    out.write(line)
                    out.write('\n')
                    # Subject
                    #for j in range(4, len(row)-1):
                    #    line = r.sub('', row[j])
                    #    out.write(line)
                    #    out.write('\n')
    finally:
        out.close()

clean_data('enron-mysqldump.sql', 'cleaned.txt')
