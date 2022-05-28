const router = require("express").Router();
const mongoose = require("mongoose");
const User = require("../models/User");

router.route("/").get((req, res) => {
  User.find()
    .then((users) => res.json(users))
    .catch((err) => res.status(400).json("Error: " + err));
});

module.exports = router;
