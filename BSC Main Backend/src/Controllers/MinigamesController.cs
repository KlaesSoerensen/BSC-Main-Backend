using BSC_Main_Backend.dto;
using Microsoft.AspNetCore.Mvc;

namespace BSC_Main_Backend.Controllers;

[ApiController]
[Route("[controller]")]
public class MinigamesController : ControllerBase
{
    private readonly ILogger<MinigamesController> _logger;
    
    public MinigamesController(ILogger<MinigamesController> logger)
    {
        _logger = logger;
    }
    
    /// <summary>
    /// Get info about some minigame alongside its difficulties.
    /// </summary>
    /// <param name="minigameId">The ID of the minigame</param>
    /// <returns></returns>
    [HttpGet("{minigameId}", Name = "GetMinigameInfo")]
    public MiniGameInfoResponseDTO getMinigameInfo([FromRoute] uint minigameId)
    {
        return null;
    }

}