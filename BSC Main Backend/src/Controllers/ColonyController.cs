using Microsoft.AspNetCore.Mvc;
using BSC_Main_Backend.dto.response;

namespace BSC_Main_Backend.Controllers;

[ApiController]
[Route("[controller]")]

public class ColonyController
{
    private readonly ILogger<ColonyController> _logger;

    public ColonyController(ILogger<ColonyController> logger)
    {
        _logger = logger;
    }


}