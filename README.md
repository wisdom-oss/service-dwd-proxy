<div align="center">
<img height="150px" src="https://raw.githubusercontent.com/wisdom-oss/brand/main/svg/standalone_color.svg">
<h1>DWD REST-Adapter</h1>
<h3>dwd-rest-proxy</h3>
<p>⚙️ A microservice for requesting data from the DWD Open data portal</p>
<img src="https://img.shields.io/github/go-mod/go-version/wisdom-oss/service-dwd-proxy?style=for-the-badge" alt="Go Lang Version"/>
<a href="openapi.yaml">
<img src="https://img.shields.io/badge/Schema%20Version-3.0.0-6BA539?style=for-the-badge&logo=OpenAPI%20Initiative" alt="Open API Schema Version"/></a>
</div>

## About
This microservice allows an easy way to access the publicly available open data
platform and the files present on it.
Since the DWD only publishes their data as compressed files and not in a 
machine-readable format (JSON or [Enterprise JSON](https://wikipedia.org/wiki/XML)).
Therefore, the service will take a request and parse the files available in the
Open Data portal and transform them into JSON.

## Credits
All datasets obtained via this microservice are provided by the DWD. The DWD
reserves all rights.

<p align="center">
<img width="200px" src="https://www.dwd.de/DE/service/copyright/dwd-logo-png.png?__blob=publicationFile&v=4">
</p>