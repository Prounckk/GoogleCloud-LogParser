#GOOGLE CLOUD FUNCTIONS
https://cloud.google.com/functions/

the function take json file with logs from google cloud bucket, scan line by line to get `TextPlayload`
Write to a txt file and send the file to FTP


##You will need: 
1) create an export of log files to google cloud bucket
see here how:  https://cloud.google.com/logging/docs/export/configure_export_v2

2) Deploy the code (without main.go! )to the Cloud Functions and set up variables

##Enthronement variables:
'FTPHOST'  
'FTPLOGIN'  
'FTPPASS'  
'FTPFOLDER'




