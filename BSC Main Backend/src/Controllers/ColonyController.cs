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

    /// <summary>
    /// Opens a colony and returns relevant details.
    /// </summary>
    /// <param name="colonyId">The ID of the colony to open</param>
    /// <param name="playerId">The ID of the player opening the colony</param>
    /// <returns>OpenColonyResponseDTO</returns>
    [HttpGet("{colonyId}/open")]
    public OpenColonyResponseDTO OpenColony([FromRoute] uint colonyId, [FromQuery] uint playerId)
    {
        return null; // Implementation here
    }

    /// <summary>
    /// Joins a colony using the provided code.
    /// </summary>
    /// <param name="code">The join code for the colony</param>
    /// <returns>JoinColonyResponseDTO</returns>
    [HttpGet("join/{code}")]
    public JoinColonyResponseDTO JoinColony([FromRoute] uint code)
    {
        return null; // Implementation here
    }
}