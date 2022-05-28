const mongoose = require("mongoose");
const bcrypt = require("bcryptjs");

const UserSchema = new mongoose.Schema(
  {
    username: {
      type: String,
      required: [true, "Please provide a username"],
      unique: true,
    },
    password: {
      type: String,
      required: [true, "Please provide a password"],
      minlength: 6,
      select: false,
    },
  },
  {
    timestamps: true,
  }
);

//Runs before the document is created
UserSchema.pre("save", async function (next) {
  //If password has been modified, it won't rehash it
  if (!this.isModified("password")) {
    next();
  }

  //Hashes the password and saves it in document
  const salt = await bcrypt.genSalt(10);
  this.password = await bcrypt.hash(this.password, salt);
  next();
});

UserSchema.methods.matchPasswords = async function (password) {
  return await bcrypt.compare(password, this.password);
};

const User = mongoose.model("User", UserSchema, "TestCollection");

module.exports = User;
