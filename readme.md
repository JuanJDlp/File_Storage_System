<center> <h1> Go - File_Storage_System </h1> </center>

<!-- ABOUT THE PROJECT -->
## About The Project

The File Storage System is an application that allows you to save files locally on your machine. It was created as a way to transfer files between devices, such as your computer and phone, or share them with family members. Additionally, it serves as a convenient way to store data that you may need later. The only limitation of this project is the available disk space on the host machine.

I intend to deploy this system on my Raspberry Pi 5 and make it accessible across the network, enabling access to my files from anywhere.

It also counts with security (JWT) and a database integration for users (Postgres)

**If you have any recommendations please let me know, i am always happy to improve**

<p align="right">(<a href="#readme-top">back to top</a>)</p>


### Built With

[![image](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![image](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- GETTING STARTED -->
## Getting Started

This is how you can set up your project locally, feel free to host it and use it however you desire.



### Prerequisites

#### Go (Golang)

Download and install Go by following the official installation guide: [Go Installation](https://go.dev/doc/install).

For Linux, you can use the following commands to download and extract the Go tarball:

```sh
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
```

#### PostgreSQL

Follow the installation guide provided by PostgreSQL: [PostgreSQL Installation Guide](https://www.postgresql.org/download/).

For setting up PostgreSQL on Ubuntu, refer to this tutorial: [How To Install and Use PostgreSQL](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-postgresql-on-ubuntu-18-04).


### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/JuanJDlp/File_Storage_System
   ```
2. Cd into the project  
    ```sh
    cd File_Storage_System
   ```
3. Create the database in postgres
   ```sh
   CREATE DATABASE file_storage;
   ```
4. Build the project
   ```sh
   go build main.go
   ```
5. Run the project 
   ```sh
   main
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

Use this project however you want; it is designed to be a Google Drive clone, allowing you to store multiple files on a machine. If you wish to create a front end for it, you are more than welcome to do so.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap

- [ ] Improve the authentication using JWT (add refresh tokens)
- [ ] Implement Cloudflare tunnels for internet access
- [ ] Deploy the system on a Raspberry Pi
   - [ ] Implement HTTPS
   - [ ] Set up Fail2Ban
- [ ] Implement multiple file downloads


<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>



