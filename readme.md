mutation {
  createUser(input: {
    name: "asad Rahman"
    age: 25
  }) {
    id
    name
    age
  }
}



{
  users {
    id
    name
    age
  }
}