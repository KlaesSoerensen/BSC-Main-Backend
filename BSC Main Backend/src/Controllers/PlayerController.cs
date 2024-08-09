using Microsoft.AspNetCore.Mvc;

namespace BSC_Main_Backend.Controllers;

[ApiController]
[Route("[controller]")]
public class PlayerController : ControllerBase
{
    
    private readonly ILogger<PlayerController> _logger;
    
    public PlayerController(ILogger<PlayerController> logger)
    {
        _logger = logger;
    }
    
    
}