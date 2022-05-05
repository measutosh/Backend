package main

import (
	"fmt"
  "os"
  "encoding/json"
  "sync"
  "io/ioutil"
  "path/filepath"
  "github.com/jcelliott/lumber"
)

const version = "1.0.0"

type (
  Logger interface {
    Fatal(string, ...interface{})
    Error(string, ...interface{})
    Warn(string, ...interface{})
    Info(string, ...interface{})
    Debug(string, ...interface{})
    Trace(string, ...interface{})
  }

  Driver struct {
    mutex   sync.Mutex
    mutexes map[string]*sync.Mutex
    dir     string
    log     Logger
  }
)

func New(dir string, options *Options)(*Driver, error) {
  
  dir = filepath.Clean(dir)

  opts := Options{}
  if options != nil {
    opts = *options
  }

  if opts.Logger == nil {
    opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
  }

  driver := Driver {
    dir :    dir,
    mutexes: make(map[string]*sync.Mutex),
    log:     opts.Logger,
  }

  if _, err := os.Stat(dir); err == nil {
    opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
    return &driver, nil
  }

  opts.Logger.Debug("Creating the database at '%s'...\n", dir)
  return &driver, os.MkdirAll(dir, 0755)
  
  
}

func (d *Driver) Write(collection, resource string, v interface{}) error {

  // the values we got in this function will be put inside a json file
  
  if collection == "" {
    return fmt.Errorf("Missing collection - no place to save record")
  }

  if resource == "" {
    return fmt.Errorf("Missing resource - unable to save record (no name)!")
  }

  mutex := d.getOrCreateMutex(collection)
  mutex.Lock()
  defer mutex.Unlock()

  // creating the json file
  dir := filepath.Join(d.dir, collection)
  fnlPath := filepath.Join(dir, resource + ".json")
  tmpPath := fnlPath + ".tmp"

  if err := os.MkdirAll(dir, 0755); err != nil {
    return err 
  }

  b, err := json.MarshalIndent(v, "", "\t")
  if err != nil {
    return err
  }

  // ensures that everything is written in next line
  b = append(b, byte('\n'))

  // writing inside the json file
  if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
    return err
  }

  return os.Rename(tmpPath, fnlPath)
}

func (d *Driver) Read(collection, resource string, v interface{}) error {

  if collection == "" {
    return fmt.Errorf("Missing collection - unable to read!")
  }

  if resource == "" {
    return fmt.Errorf("Missing resource - unable to save record(no name)!")
  }

  // to get the record
  record := filepath.Join(d.dir, collection, resource)

  // if record doesn't exist send error or read the record
  if _, err := stat(record); err != nil {
    return err
  }

  // capture the value at b
  b, err := ioutil.ReadFile(record + ".json")
  if err != nil {
    return err
  }

  return json.Unmarshal(b, &v)
}

func (d *Driver) ReadAll(collection string)([]string, error) {
  // the slice string is all the data which will be returned

  if collection == ""{
    return nil, fmt.Errorf("Missing collection - unable to read")
  }

  // enters into the directory and gets the collection name
  dir := filepath.Join(d.dir, collection)

  if _, err := stat(dir); err != nil {
    return nil, err
  }

  // reading all the records, the entire directory
  // cathcing values in files and error in blank itor
  files, _ := ioutil.ReadDir(dir)

  var records []string

  // getting all names then reading them
  for _, file := range files {
    b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
    if err != nil {
      return nil, err
    }

    records = append(records, string(b))
  }

  return records, nil
}

func (d *Driver) Delete(collection, resource string) error {
  
  path := filepath.Join(collection, resource)
  mutex := d.getOrCreateMutex(collection)
  mutex.Lock()
  defer mutex.Unlock()

  dir := filepath.Join(d.dir, path)

  switch fi, err := stat(dir); {
    case fi == nil, err != nil:
      return fmt.Errorf("Unable to find file or directory named %v\n", path)
    
    case fi.Mode().IsDir():
      return os.RemoveAll(dir)
    
    case fi.Mode().IsRegular():
      return os.RemoveAll(dir + ".json")
  }

  return nil  
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {

  d.mutex.Lock()
  defer d.mutex.Unlock()
  m, ok := d.mutexes[collection]

  // if mutex doesn't exist then create a new mutex
  if !ok {
    m = &sync.Mutex{}
    d.mutexes[collection] = m
  }

  return m
}

func stat(path string)(fi os.FileInfo, err error) {
  
  if fi, err = os.Stat(path); os.IsNotExist(err) {
    fi, err = os.Stat(path + ".json")
  }
  return
}

type Options struct {
  Logger
}

type Address struct {
  City    string
  State   string
  Country string
  Pincode json.Number
}

type User struct {
  Name    string
  Age     json.Number
  Contact string
  Company string
  Address Address
}

func main() {
  dir := "./"

  db, err := New(dir, nil)
  if err != nil {
    	fmt.Println("Error", err)
  }

  employees := []User{
    {"Rocky", "23", "9238478234", "Zomato", Address{"Bangalore", "Karnataka", "India", "823478"}},
    {"Chitti", "25", "900078234", "Swiggy", Address{"Chennai", "Tamilnadu", "India", "823478"}},
    {"Ajay", "32", "9238999994", "Byjus", Address{"Mumbai", "Maharstra", "India", "823478"}},
    {"Bunny", "34", "7778478234", "Razorpay", Address{"Delhi", "Newdelhi", "India", "823478"}},
    {"Arya", "43", "9238474444", "Dukaan", Address{"Kolkata", "Westbengal", "India", "823478"}},
    {"Shyam", "27", "9000000004", "Flipkart", Address{"Lucknow", "Utterpradesh", "India", "823478"}},
  }

  for _, value := range employees {
    db.Write("users", value.Name, User{
      Name:    value.Name,
      Age:     value.Age,
      Contact: value.Contact,
      Company: value.Company,
      Address: value.Address,
    })
  }

  records, err := db.ReadAll("users")
  if err != nil {
    fmt.Println("Error", err)
  }
  fmt.Println(records)

  allusers := []User{}

  for _, f := range records {
    employeeFound := User{}
    if err := json.Unmarshal([]byte(f), &employeeFound); err != nil{
      fmt.Println("Error", err)
    }

    allusers = append(allusers, employeeFound)
  }

  fmt.Println(allusers)


  // // to delete one user
  // if err := db.Delete("user", "Rocky"); err != nil {
  //   fmt.Println("Error", err)
  // }

  // // to delete all users
  // if err := db.Delete("user", ""); err != nil {
  //   fmt.Println("Error", err)
  // }
  
}
