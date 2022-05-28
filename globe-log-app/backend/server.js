const express = require("express");
const cors = require("cors");
const mongoose = require("mongoose");

require("dotenv").config();

const app = express();
const port = process.env.PORT || 5000;

app.use(cors());
app.use(express.json());

// Create connection to mongodb

const uri = process.env.ATLAS_URI;
mongoose.connect(uri, {
  useNewUrlParser: true,
  useUnifiedTopology: true,
});
const connection = mongoose.connection; // Instance of db connection
connection.once("open", () => {
  console.log("MongoDB database connection established successfully");
});

app.use("/auth", require("./routes/auth"));
app.use("/users", require("./routes/users"));

app.listen(port, () => {
  console.log(`Server is running on port:${port}`);
});
