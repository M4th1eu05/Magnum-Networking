for _, file in ipairs(os.files("*.c")) do
    target(path.basename(file), function ()
        set_group("examples")
        add_files(file)
        add_deps("enet6")
    end)
end
