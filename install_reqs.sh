set -e
rm -r ./cas
python3 -m venv cas
source ./cas/bin/activate
pip install --upgrade pip
pip install -r requirements.txt