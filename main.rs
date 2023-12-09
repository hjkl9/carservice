use std::io;
use std::File;

fn main() {
    // initialization
}

fn read_config_from_file() -> Result<String, io::Error> {
    let mut s = String::new();
    File::open("./etc/carservice.yaml")?.read_to_string(&mut s)?;
    Ok(s)
}