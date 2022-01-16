package main

import (
         "bytes"
         "log"
         "fmt"
         "time"
         "golang.org/x/crypto/ssh"
         "github.com/360EntSecGroup-Skylar/excelize"
) 

func main() {

     f, err := excelize.OpenFile("ntp_go.xlsx")  // open Excel file
     if err != nil {
        panic(err)
     }

     sheet := f.GetSheetName(0)            // Get the sheet
     rows, err := f.GetRows(sheet)         // get the number of row
     if err != nil {
        fmt.Println(err)
     }

     num :=  len(rows)  // get the number of rows accordingly,  loop will execute

     for i := 2; i <= num; i++ {
         ip, _  := f.GetCellValue(sheet, fmt.Sprintf("B%d", i))  // get IP address
         user, _ := f.GetCellValue(sheet, fmt.Sprintf("C%d", i))  // get user name
         password, _ := f.GetCellValue(sheet, fmt.Sprintf("D%d", i))  // get password 
         fmt.Printf("%s\n", ip)
     
        status, offset, err :=  checkNTP(ip, user, password)    // Calling the function to get status of NTP Status and Offset value
        if err != nil {
        fmt.Println("Got Some Error:", err.Error())
        }

        fmt.Println(status)  // print on cosole
        fmt.Println(offset)  // print on console
 
        f.SetCellValue(sheet, fmt.Sprintf("E%d", i), status)  // set in excel sheet
        f.SetCellValue(sheet, fmt.Sprintf("F%d", i), offset)  // set in excel sheet
        f.SetCellValue(sheet, fmt.Sprintf("G%d", i), err   )  // in comment field error for particular server
    }    

    err = f.SaveAs("./ntp_go.xlsx")  // save the excel with data
     if err != nil {
        fmt.Println(err)
     }
 

}

func checkNTP(ip string, user string, password string) (status string, offset string, err error) { 
     config := &ssh.ClientConfig {
          User: user, 
          Auth: []ssh.AuthMethod { 
                ssh.Password(password), 
          },
          HostKeyCallback: ssh.InsecureIgnoreHostKey(), // it accepts every hosts
          Timeout: 5 * time.Second, 
     }

     client, err := ssh.Dial("tcp", ip+":22", config) 
     if err != nil {
//        panic("Failed to Dial: "+ err.Error())
       return "", "", err
     }

     session1, err := client.NewSession()
     if err != nil {
       return "", "", err
     }

     defer session1.Close()

     var b bytes.Buffer
     session1.Stdout = &b
     if err := session1.Run("timedatectl  | grep -i sync"); err != nil {
        log.Println("Failed to run on client: " + err.Error())
        return "", "", err
     }


      session2, err := client.NewSession()
      if err != nil {
        return "", "",  err
      }

     var c bytes.Buffer
     session2.Stdout = &c
     if err := session2.Run("/usr/sbin/ntpq -p | awk '{print $9}' | awk 'FNR == 3 {print}'"); err != nil {
        log.Println("Failed to run on client: " + err.Error())
        return  "", "", err
     }

     return b.String(), c.String(), nil
}

