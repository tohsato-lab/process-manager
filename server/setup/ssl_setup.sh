openssl genrsa 2048 > key.pem
# openssl req -new -key key.pem > csr.pem
openssl req -new -subj '/' -key key.pem -out csr.pem
openssl x509 -days 3650 -req -signkey key.pem < csr.pem > cert.pem
rm csr.pem
