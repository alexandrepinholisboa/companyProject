# Data integration project

This project was created to handle the necessity of the client to combine data from different sources.

## The main functionality

Must be possible to upload a CSV file with the headers "name" and "addresszip" successfully.
After that, a new CSV file shall be merged, containing headers "name", "addresszip" and "website".
The "website' value will be added within the existing records when the "name" and "addresszip" values exactly match with the DB.

## Notes

- The 'companyProject' shall be cloned within '$GOPATH/src'
- Before start the application, run the test case contained within 'src/companyProject/tests'. Thereby, verifying if the API's are working correctly.
- The file "Makefile" contains the commands to start/test the API
- The folder "files" contains the files to be imported
    - File to be upload:
        - The file 'q1_catalog.csv' is the file to be upload that contains the headers 'name' and 'addresszip'
    - File to be merged:
        - The file 'q2_clientData.csv' is an invalid file that contains the header 'name', 'addresszip' and 'website'
        - The file 'q3_clientData.csv' is a valid file that contains the header 'name', 'zip' and 'website'
        - The file 'q4_clientData.txt' is an invalid txt file
- The content to be merged must match some criterias. Otherwise, it will be skiped
    - The file must be CSV
    - The value within the 'name' must be in upper case
    - The value within the 'zip' must have 5 digits
    - The value within the 'website' must be in lower case

