counter = 0

request = function()
    counter = counter + 1
    local email = "userm" .. counter .. "@maildrop.cc"
    local body = string.format('{"email":"%s","password":"password123"}', email)

    return wrk.format("POST", "/api/auth/register", {["Content-Type"] = "application/json"}, body)
end
