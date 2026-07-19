export function validateRegister(data) {

    if (data.gender !== "Male" && data.gender !== "Female") {
        return "Invalid gender. Please select 'Male' or 'Female'";
    }

    if (
    !data.nickname ||
    !data.first_name ||
    !data.last_name ||
    !data.email ||
    !data.password ||
    !data.confirmpassword
    ) {
    return "Please fill all the fields";
    }
    if(data.age < 9 || data.age > 150){
        return "age must be between 9 and 150"
    }

    if (
        data.nickname.length < 3 || data.nickname.length > 50 ||
        data.first_name.length < 3 || data.first_name.length > 50 ||
        data.last_name.length < 3 || data.last_name.length > 50
    ) {
        return "Nickname, first name, and last name must be between 3 and 50 characters long";
    }

    if (data.password.length < 8 || data.password.length > 50) {
        return "Password must be between 8 and 50 characters long";
    }
    if (data.password !== data.confirmpassword) {
    return "Passwords do not match";
    }

    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

    if (!emailRegex.test(data.email)) {
        return "Invalid email address";
    }

    return null;
}